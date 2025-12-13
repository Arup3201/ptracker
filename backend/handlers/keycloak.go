package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ptracker/apierr"
	"github.com/ptracker/db"
	"github.com/ptracker/utils"
)

const (
	KC_URL           = "http://localhost:8080"
	KC_REALM         = "ptracker"
	KC_CLIENT_ID     = "api"
	KC_CLIENT_SECRET = "cp50avHQeX18cESEraheJvr3RhUBMq2A"
	KC_REDIRECT_URI  = "http://localhost:8081/api/keycloak/callback"
)

type ApiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ApiData map[string]any

type ApiResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	ApiError `json:"error,omitempty"`
	ApiData  `json:"data,omitempty"`
}

type KCError struct {
	ErrorCode        string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

var states = map[string]string{}

func KeycloakLogin(w http.ResponseWriter, r *http.Request) {
	verifier, _ := utils.CreateCodeVerifier()
	challenge := verifier.CodeChallengeSHA256()

	state := uuid.NewString()
	states[state] = verifier.Value

	kc_auth_url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/auth?"+
		"scope=openid"+
		"&response_type=code"+
		"&client_id=%s"+
		"&redirect_uri=%s"+
		"&state=%s"+
		"&code_challenge_method=S256"+
		"&code_challenge=%s",
		KC_URL, KC_REALM, KC_CLIENT_ID, KC_REDIRECT_URI, state, challenge)
	http.Redirect(w, r, kc_auth_url, http.StatusSeeOther)
}

func KeycloakCallback(w http.ResponseWriter, r *http.Request) {
	authorization_code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	authorization_error_code := r.URL.Query().Get("error")

	if authorization_error_code != "" {
		authorization_error_description := r.URL.Query().Get("error_description")

		fmt.Printf("[ERROR] authorization code request error: %s", authorization_error_description)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ApiResponse{
			Status:  "Error",
			Message: "Keycloak authorization error",
			ApiError: ApiError{
				Code:    authorization_error_code,
				Message: authorization_error_description,
			},
		})
		return
	}

	tokenEndpoint := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", KC_URL, KC_REALM)
	if _, ok := states[state]; !ok {
		fmt.Printf("[ERROR] PKCE code verifier missing")

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ApiResponse{
			Status:  "Error",
			Message: "Authentication error",
			ApiError: ApiError{
				Code:    "malicious_attempt",
				Message: "No code_verifier found for the state",
			},
		})
		return
	}
	res, err := http.PostForm(tokenEndpoint, url.Values{
		"grant_type":    []string{"authorization_code"},
		"code":          []string{authorization_code},
		"code_verifier": []string{states[state]},
		"redirect_uri":  []string{KC_REDIRECT_URI},
		"client_id":     []string{KC_CLIENT_ID},
		"client_secret": []string{KC_CLIENT_SECRET},
	})
	delete(states, state)

	if err != nil {
		fmt.Printf("[ERROR] token request error: %s", err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ApiResponse{
			Status:  "Error",
			Message: "Keycloak token endpoint error",
			ApiError: ApiError{
				Code:    "internal_error",
				Message: "Error requesting token from keycloak",
			},
		})
		return
	}

	if res.StatusCode != http.StatusOK {
		var KCErrorResponse KCError
		json.NewDecoder(res.Body).Decode(&KCErrorResponse)

		fmt.Printf("[ERROR] %s", KCErrorResponse.ErrorDescription)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ApiResponse{
			Status:  "Error",
			Message: "Keycloak token endpoint error",
			ApiError: ApiError{
				Code:    KCErrorResponse.ErrorCode,
				Message: "Keycloak error response for token",
			},
		})
		return
	}

	var TokenResponse struct {
		AccessToken      string `json:"access_token"`
		TokenType        string `json:"token_type"`
		RefreshToken     string `json:"refresh_token"`
		RefreshExpiresIn int    `json:"refresh_expires_in"`
		IDToken          string `json:"id_token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&TokenResponse); err != nil {
		fmt.Printf("[ERROR] token unpacking error: %s", err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ApiResponse{
			Status:  "Error",
			Message: "Keycloak token endpoint error",
			ApiError: ApiError{
				Code:    "internal_error",
				Message: "Error unpacking token payload",
			},
		})
		return
	}

	userinfoEndpoint := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/userinfo", KC_URL, KC_REALM)
	req, _ := http.NewRequest("GET", userinfoEndpoint, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", TokenResponse.AccessToken))
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("[ERROR] userinfo request error: %s", err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ApiResponse{
			Status:  "Error",
			Message: "Keycloak userinfo endpoint error",
			ApiError: ApiError{
				Code:    "internal_error",
				Message: "Error requesting userinfo endpoint",
			},
		})
		return
	}

	if res.StatusCode != http.StatusOK {
		var KCErrorResponse KCError
		json.NewDecoder(res.Body).Decode(&KCErrorResponse)

		fmt.Printf("[ERROR] status: %d: %s", res.StatusCode, KCErrorResponse.ErrorDescription)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ApiResponse{
			Status:  "Error",
			Message: "Keycloak userinfo endpoint error",
			ApiError: ApiError{
				Code:    KCErrorResponse.ErrorCode,
				Message: "Keycloak error response for userinfo",
			},
		})
		return
	}

	var keycloakUserInfo struct {
		Subject   string `json:"sub"`
		Name      string `json:"name"`
		Username  string `json:"preferred_username"`
		Email     string `json:"email"`
		AvatarUrl string `json:"picture"`
	}
	if err := json.NewDecoder(res.Body).Decode(&keycloakUserInfo); err != nil {
		fmt.Printf("[ERROR] userinfo unpacking error: %s", err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ApiResponse{
			Status:  "Error",
			Message: "Keycloak userinfo unpacking error",
			ApiError: ApiError{
				Code:    "internal_error",
				Message: "Error unpacking userinfo response",
			},
		})
		return
	}

	user, err := db.FindUserWithIdp(keycloakUserInfo.Subject, "keycloak")
	if errors.Is(err, &apierr.ResourceNotFound{}) {
		user, err = db.CreateUser(keycloakUserInfo.Subject, "keycloak",
			keycloakUserInfo.Name, keycloakUserInfo.Name, keycloakUserInfo.Email,
			keycloakUserInfo.AvatarUrl)
		if err != nil {
			fmt.Printf("[ERROR] new user create error: %s", err)

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ApiResponse{
				Status:  "Error",
				Message: "User create error",
				ApiError: ApiError{
					Code:    "internal_error",
					Message: "Tried creating a new user, but failed",
				},
			})
			return
		}
	}

	expiresAt := time.Now().Add(time.Duration(TokenResponse.RefreshExpiresIn * int(time.Second)))
	sessionId, err := db.CreateSession(user.Id, TokenResponse.RefreshToken, r.UserAgent(),
		strings.Split(r.RemoteAddr, ":")[0], "None", expiresAt)
	if err != nil {
		fmt.Printf("[ERROR] session create error: %s", err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ApiResponse{
			Status:  "Error",
			Message: "Session create error",
			ApiError: ApiError{
				Code:    "internal_error",
				Message: "Failed to create session for the user",
			},
		})
		return
	}

	cookie := &http.Cookie{
		Name:     "PTRACKER_SESSION_ID",
		Value:    sessionId,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		Expires:  expiresAt,
	}
	http.SetCookie(w, cookie)

	json.NewEncoder(w).Encode(ApiResponse{
		Status:  "Success",
		Message: "Login success",
	})
}
