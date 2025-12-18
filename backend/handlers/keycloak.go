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

var (
	KC_URL              string
	KC_REALM            string
	KC_CLIENT_ID        string
	KC_CLIENT_SECRET    string
	KC_REDIRECT_URI     string
	ENCRYPTION_SECRET   string
	HOME_URL            string
	SESSION_COOKIE_NAME = "PTRACKER_SESSION_ID"
)

type KCError struct {
	ErrorCode        string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

var states = map[string]string{}

func KeycloakLogin(w http.ResponseWriter, r *http.Request) error {
	verifier, _ := utils.CreateCodeVerifier()
	challenge := verifier.CodeChallengeSHA256()

	state := uuid.NewString()
	states[state] = verifier.Value

	if KC_URL == "" || KC_REALM == "" || KC_CLIENT_ID == "" || KC_REDIRECT_URI == "" {
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Server encountered an issue",
			Err:     fmt.Errorf("keycloak login: one or more of the required env KC_URL, KC_REALM, KC_CLIENT_ID and KC_REDIRECT_URI missing"),
		}
	}

	if state == "" || challenge == "" {
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Server encountered an issue",
			Err:     fmt.Errorf("keycloak login: state/challenge is missing"),
		}
	}

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

	return nil
}

func KeycloakCallback(w http.ResponseWriter, r *http.Request) error {
	if KC_URL == "" || KC_REALM == "" || KC_CLIENT_ID == "" || KC_REDIRECT_URI == "" || KC_CLIENT_SECRET == "" || HOME_URL == "" {
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Server encountered an issue",
			Err:     fmt.Errorf("keycloak login: one or more of the required env KC_URL, KC_REALM, KC_CLIENT_ID, KC_REDIRECT_URI and HOME_URL missing"),
		}
	}

	if ENCRYPTION_SECRET == "" {
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Server encountered an issue",
			Err:     fmt.Errorf("keycloak login: ENCRYPTION_SECRET env missing"),
		}
	}

	authorization_code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	authorization_error_code := r.URL.Query().Get("error")

	if authorization_error_code != "" {
		authorization_error_description := r.URL.Query().Get("error_description")
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Authorization denied",
			Err:     fmt.Errorf("keycloak callback: %s", authorization_error_description),
		}
	}

	if _, ok := states[state]; !ok {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Authorization denied for repeated PKCE",
			Err:     fmt.Errorf("keycloak callback: code_verifier state reused"),
		}
	}

	tokenEndpoint := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", KC_URL, KC_REALM)
	res, err := http.PostForm(tokenEndpoint, url.Values{
		"grant_type":    []string{"authorization_code"},
		"code":          []string{authorization_code},
		"code_verifier": []string{states[state]},
		"redirect_uri":  []string{KC_REDIRECT_URI},
		"client_id":     []string{KC_CLIENT_ID},
		"client_secret": []string{KC_CLIENT_SECRET},
	})
	if err != nil {
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Authorization failed",
			Err:     fmt.Errorf("keycloak callback: keycloak token request: %w", err),
		}
	}
	defer res.Body.Close()
	delete(states, state)

	if res.StatusCode != http.StatusOK {
		var KCErrorResponse KCError
		if err := json.NewDecoder(res.Body).Decode(&KCErrorResponse); err != nil {
			return &HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Authorization failed",
				Err:     fmt.Errorf("keycloak callback: keycloak token response error decode: %w", err),
			}
		}

		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Authorization failed",
			Err:     fmt.Errorf("keycloak callback: keycloak token response: %s", KCErrorResponse.ErrorDescription),
		}
	}

	var TokenResponse struct {
		AccessToken      string `json:"access_token"`
		TokenType        string `json:"token_type"`
		RefreshToken     string `json:"refresh_token"`
		RefreshExpiresIn int    `json:"refresh_expires_in"`
		IDToken          string `json:"id_token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&TokenResponse); err != nil {
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Authorization failed",
			Err:     fmt.Errorf("keycloak callback: keycloak token response body decode: %w", err),
		}
	}

	userinfoEndpoint := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/userinfo", KC_URL, KC_REALM)
	req, _ := http.NewRequest("GET", userinfoEndpoint, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", TokenResponse.AccessToken))
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Authorization failed",
			Err:     fmt.Errorf("keycloak callback: keycloak userinfo request: %w", err),
		}
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var KCErrorResponse KCError
		if err := json.NewDecoder(res.Body).Decode(&KCErrorResponse); err != nil {
			return &HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Authorization failed",
				Err:     fmt.Errorf("keycloak callback: keycloak userinfo response error decode: %w", err),
			}
		}

		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Authorization failed",
			Err:     fmt.Errorf("keycloak callback: keycloak userinfo response: %s", KCErrorResponse.ErrorDescription),
		}
	}

	var keycloakUserInfo struct {
		Subject   string `json:"sub"`
		Name      string `json:"name"`
		Username  string `json:"preferred_username"`
		Email     string `json:"email"`
		AvatarUrl string `json:"picture"`
	}
	if err := json.NewDecoder(res.Body).Decode(&keycloakUserInfo); err != nil {
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Authorization failed",
			Err:     fmt.Errorf("keycloak callback: keycloak userinfo response body decode: %w", err),
		}
	}

	user, err := db.FindUserWithIdp(keycloakUserInfo.Subject, "keycloak")
	if errors.Is(err, apierr.ErrResourceNotFound) {
		user, err = db.CreateUser(keycloakUserInfo.Subject, "keycloak",
			keycloakUserInfo.Name, keycloakUserInfo.Name, keycloakUserInfo.Email,
			keycloakUserInfo.AvatarUrl)
		if err != nil {
			return &HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Authorization failed",
				Err:     fmt.Errorf("keycloak callback: new user create: %w", err),
			}
		}
	}

	tokenExpiresAt := time.Now().Add(time.Duration(TokenResponse.RefreshExpiresIn * int(time.Second)))

	encryptedRefreshToken, err := utils.EncryptAES([]byte(TokenResponse.RefreshToken), []byte(ENCRYPTION_SECRET))
	if err != nil {
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Authorization failed",
			Err:     fmt.Errorf("keycloak callback: refresh token encryption: %w", err),
		}
	}
	session, err := db.CreateSession(user.Id, encryptedRefreshToken, r.UserAgent(),
		strings.Split(r.RemoteAddr, ":")[0], "None", tokenExpiresAt)
	if err != nil {
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Authorization failed",
			Err:     fmt.Errorf("keycloak callback: new session create: %w", err),
		}
	}

	accessTokens[session.Id] = TokenResponse.AccessToken

	cookie := &http.Cookie{
		Name:     SESSION_COOKIE_NAME,
		Value:    session.Id,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		Expires:  tokenExpiresAt,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, HOME_URL, http.StatusSeeOther)
	return nil
}

func KeycloakRefresh(w http.ResponseWriter, r *http.Request) error {
	cookies := r.Cookies()
	sessionId, err := utils.GetSessionIdFromCookie(cookies, SESSION_COOKIE_NAME)
	if err != nil {
		return &HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "User session has expired",
			Err:     fmt.Errorf("keycloak refresh token: %w", err),
		}
	}
	session, err := db.GetActiveSession(sessionId)
	if err != nil {
		return &HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "User session has expired",
			Err:     fmt.Errorf("keycloak refresh token: %w", err),
		}
	}

	refreshToken, err := utils.DecryptAES(session.RefreshTokenEncrypted, []byte(ENCRYPTION_SECRET))
	if err != nil {
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Token refresh error",
			Err:     fmt.Errorf("keycloak refresh token: %w", err),
		}
	}

	tokenEndpoint := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", KC_URL, KC_REALM)
	res, err := http.PostForm(tokenEndpoint, url.Values{
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{string(refreshToken)},
		"redirect_uri":  []string{KC_REDIRECT_URI},
		"client_id":     []string{KC_CLIENT_ID},
		"client_secret": []string{KC_CLIENT_SECRET},
		"scope":         []string{"openid"},
	})
	if err != nil {
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Token refresh error",
			Err:     fmt.Errorf("keycloak refresh token: keycloak refresh token request: %w", err),
		}
	}

	if res.StatusCode != http.StatusOK {
		// revoke session and remove from cookie

		err := db.MakeSessionInactive(session.Id)
		if err != nil {
			return &HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Token refresh error",
				Err:     fmt.Errorf("keycloak refresh token: %w", err),
			}
		}

		cookie := &http.Cookie{
			Name:     SESSION_COOKIE_NAME,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteDefaultMode,
			Expires:  time.Unix(0, 0),
		}
		http.SetCookie(w, cookie)

		var KCErrorResponse KCError
		if err := json.NewDecoder(res.Body).Decode(&KCErrorResponse); err != nil {
			return &HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "User session expired",
				Err:     fmt.Errorf("keycloak refresh token: keycloak error response: %w", err),
			}
		}

		return &HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "User session expired",
			Err:     fmt.Errorf("keycloak refresh token: %w", err),
		}
	}

	var TokenResponse struct {
		AccessToken      string `json:"access_token"`
		TokenType        string `json:"token_type"`
		RefreshToken     string `json:"refresh_token"`
		RefreshExpiresIn int    `json:"refresh_expires_in"`
		IDToken          string `json:"id_token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&TokenResponse); err != nil {
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Token refresh error",
			Err:     fmt.Errorf("keycloak refresh token: keycloak token response: %w", err),
		}
	}

	tokenExpiresAt := time.Now().Add(time.Duration(TokenResponse.RefreshExpiresIn * int(time.Second)))
	encryptedRefreshToken, err := utils.EncryptAES([]byte(TokenResponse.RefreshToken), []byte(ENCRYPTION_SECRET))
	if err != nil {
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Token refresh error",
			Err:     fmt.Errorf("keycloak refresh token: %w", err),
		}
	}

	err = db.UpdateSession(session.Id, encryptedRefreshToken, tokenExpiresAt)
	if err != nil {
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Token refresh error",
			Err:     fmt.Errorf("keycloak refresh token: %w", err),
		}
	}
	accessTokens[session.Id] = TokenResponse.AccessToken

	json.NewEncoder(w).Encode(HTTPSuccessResponse{
		Status:  "success",
		Message: "Access token refreshed",
	})
	return nil
}

func KeycloakLogout(w http.ResponseWriter, r *http.Request) error {
	sessionId, err := utils.GetSessionIdFromCookie(r.Cookies(), SESSION_COOKIE_NAME)
	if err != nil {
		return &HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "User session has expired",
			Err:     fmt.Errorf("keycloak refresh token: %w", err),
		}
	}

	// revoke db session, remove cookie and in-memory access token

	err = db.MakeSessionInactive(sessionId)
	if err != nil {
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Logout error",
			Err:     fmt.Errorf("keycloak logout: %w", err),
		}
	}

	cookie := &http.Cookie{
		Name:     SESSION_COOKIE_NAME,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		Expires:  time.Unix(0, 0),
	}
	http.SetCookie(w, cookie)

	delete(accessTokens, sessionId)

	json.NewEncoder(w).Encode(HTTPSuccessResponse{
		Status:  "success",
		Message: "Logout success",
	})
	return nil
}
