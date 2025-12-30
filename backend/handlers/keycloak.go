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
	SESSION_COOKIE_NAME = "PTRACKER_SESSION_ID"
)

type KCToken struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	IDToken          string `json:"id_token"`
}

type KCUserInfo struct {
	Subject   string `json:"sub"`
	Name      string `json:"name"`
	Username  string `json:"preferred_username"`
	Email     string `json:"email"`
	AvatarUrl string `json:"picture"`
}

type KCError struct {
	ErrorCode        string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func GetToken(url, realm string, kcRequestValue url.Values) (*KCToken, error) {
	tokenEndpoint := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", url, realm)
	res, err := http.PostForm(tokenEndpoint, kcRequestValue)
	if err != nil {
		return nil, fmt.Errorf("keycloak token request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var KCErrorResponse KCError
		if err := json.NewDecoder(res.Body).Decode(&KCErrorResponse); err != nil {
			return nil, fmt.Errorf("keycloak token error status %d", res.StatusCode)
		}

		return nil, fmt.Errorf("keycloak token response: %s", KCErrorResponse.ErrorDescription)
	}

	var tokenResponse KCToken
	if err := json.NewDecoder(res.Body).Decode(&tokenResponse); err != nil {
		return nil, fmt.Errorf("keycloak token response body decode: %w", err)
	}

	return &tokenResponse, nil
}

func GetUserInfo(url, realm, accessToken string) (*KCUserInfo, error) {
	userinfoEndpoint := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/userinfo", url, realm)
	req, err := http.NewRequest("GET", userinfoEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("userinfo request generate: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("keycloak userinfo response: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var KCErrorResponse KCError
		if err := json.NewDecoder(res.Body).Decode(&KCErrorResponse); err != nil {
			return nil, fmt.Errorf("keycloak userinfo error status %d", res.StatusCode)
		}

		return nil, fmt.Errorf("keycloak userinfo response: %s", KCErrorResponse.ErrorDescription)
	}

	var keycloakUserInfo KCUserInfo
	if err := json.NewDecoder(res.Body).Decode(&keycloakUserInfo); err != nil {
		return nil, fmt.Errorf("keycloak userinfo response body decode: %w", err)
	}

	return &keycloakUserInfo, nil
}

func GetSessionCookie(refreshTokenExpires int,
	accessToken, refreshToken string,
	userId, userAgent, ipAddress, device string,
	encryptionKey string) (*http.Cookie, error) {
	tokenExpiresAt := time.Now().Add(time.Duration(refreshTokenExpires) * time.Second)

	encryptedRefreshToken, err := utils.EncryptAES([]byte(refreshToken), []byte(encryptionKey))
	if err != nil {
		return nil, fmt.Errorf("refresh token encryption: %w", err)
	}
	session, err := db.CreateSession(userId, encryptedRefreshToken, userAgent,
		ipAddress, device, tokenExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("new session create: %w", err)
	}

	accessTokens[session.Id] = accessToken

	cookie := &http.Cookie{
		Name:     SESSION_COOKIE_NAME,
		Value:    session.Id,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		Expires:  tokenExpiresAt,
	}

	return cookie, nil
}

type KeycloakHandler struct {
	KCUrl          string
	KCRealm        string
	KCClientId     string
	KCClientSecret string
	KCRedirectURI  string
	HomeURL        string
	EncryptionKey  string
}

func CreateKeycloakHandler(kcUrl, kcRealm, kcClientId, kcClientSecret, kcRedirectURI string,
	homeUrl, encryptionKey string) (*KeycloakHandler, error) {
	if kcUrl == "" || kcRealm == "" || kcClientId == "" || kcRedirectURI == "" {
		return nil, fmt.Errorf("none of the parameters can be empty")
	}

	return &KeycloakHandler{
		KCUrl:          kcUrl,
		KCRealm:        kcRealm,
		KCClientId:     kcClientId,
		KCClientSecret: kcClientSecret,
		KCRedirectURI:  kcRedirectURI,
		HomeURL:        homeUrl,
		EncryptionKey:  encryptionKey,
	}, nil
}

var states = map[string]string{}

func (handler *KeycloakHandler) KeycloakLogin(w http.ResponseWriter, r *http.Request) error {
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
		handler.KCUrl, handler.KCRealm, handler.KCClientId, handler.KCRedirectURI, state, challenge)
	http.Redirect(w, r, kc_auth_url, http.StatusSeeOther)

	return nil
}

func (handler *KeycloakHandler) KeycloakCallback(w http.ResponseWriter, r *http.Request) error {
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
			Err:     fmt.Errorf("keycloak callback: code_verifier state not found"),
		}
	}

	tokenResponse, err := GetToken(handler.KCUrl, handler.KCRealm, url.Values{
		"grant_type":    []string{"authorization_code"},
		"code":          []string{authorization_code},
		"code_verifier": []string{states[state]},
		"redirect_uri":  []string{handler.KCRedirectURI},
		"client_id":     []string{handler.KCClientId},
		"client_secret": []string{handler.KCClientSecret},
	})
	if err != nil {
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Authorization failed",
			Err:     fmt.Errorf("keycloak callback: %w", err),
		}
	}
	delete(states, state)

	kcUserInfo, err := GetUserInfo(handler.KCUrl, handler.KCRealm, tokenResponse.AccessToken)
	if err != nil {
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Authorization failed",
			Err:     fmt.Errorf("keycloak callback: %w", err),
		}
	}

	user, err := db.FindUserWithIdp(kcUserInfo.Subject, "keycloak")
	if errors.Is(err, apierr.ErrResourceNotFound) {
		user, err = db.CreateUser(kcUserInfo.Subject, "keycloak",
			kcUserInfo.Name, kcUserInfo.Name, kcUserInfo.Email,
			kcUserInfo.AvatarUrl)
		if err != nil {
			return &HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Authorization failed",
				Err:     fmt.Errorf("keycloak callback: new user create: %w", err),
			}
		}
	}

	cookie, err := GetSessionCookie(tokenResponse.RefreshExpiresIn, tokenResponse.AccessToken, tokenResponse.RefreshToken, user.Id, r.UserAgent(), strings.Split(r.RemoteAddr, " ")[0], "None", handler.EncryptionKey)
	if err != nil {
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Authorization failed",
			Err:     fmt.Errorf("keycloak callback: %w", err),
		}
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, handler.HomeURL, http.StatusSeeOther)
	return nil
}

func (handler *KeycloakHandler) KeycloakRefresh(w http.ResponseWriter, r *http.Request) error {
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

	refreshToken, err := utils.DecryptAES(session.RefreshTokenEncrypted, []byte(handler.EncryptionKey))
	if err != nil {
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Token refresh error",
			Err:     fmt.Errorf("keycloak refresh token: %w", err),
		}
	}

	tokenEndpoint := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", handler.KCUrl, handler.KCRealm)
	res, err := http.PostForm(tokenEndpoint, url.Values{
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{string(refreshToken)},
		"redirect_uri":  []string{handler.KCRedirectURI},
		"client_id":     []string{handler.KCClientId},
		"client_secret": []string{handler.KCClientSecret},
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
				Err:     fmt.Errorf("keycloak refresh token: response error status: %d", res.StatusCode),
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
	encryptedRefreshToken, err := utils.EncryptAES([]byte(TokenResponse.RefreshToken), []byte(handler.EncryptionKey))
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

func (handler *KeycloakHandler) KeycloakLogout(w http.ResponseWriter, r *http.Request) error {
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
