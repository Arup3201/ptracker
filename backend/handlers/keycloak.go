package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
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
		var KCErrorResponse struct {
			ErrorCode        string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		json.NewDecoder(res.Body).Decode(&KCErrorResponse)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ApiResponse{
			Status:  "Error",
			Message: "Keycloak token endpoint error",
			ApiError: ApiError{
				Code:    KCErrorResponse.ErrorCode,
				Message: "Keycloak error response for token",
			},
		})
		fmt.Printf("[ERROR] %s", KCErrorResponse.ErrorDescription)
		return
	}

	var TokenResponse struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		IDToken      string `json:"id_token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&TokenResponse); err != nil {
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

	json.NewEncoder(w).Encode(ApiResponse{
		Status:  "Success",
		Message: "Login success",
		ApiData: ApiData{
			"access_token":  TokenResponse.AccessToken,
			"refresh_token": TokenResponse.RefreshToken,
			"token_type":    TokenResponse.TokenType,
			"expires_in":    TokenResponse.ExpiresIn,
			"id_token":      TokenResponse.IDToken,
		},
	})
}
