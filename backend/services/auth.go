package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/ptracker/apierr"
	"github.com/ptracker/domain"
	"github.com/ptracker/interfaces"
	"github.com/ptracker/utils"
)

const (
	DEFAULT_PROVIDER = "keycloak"
)

type Token struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	IDToken          string `json:"id_token"`
}

type UserInfo struct {
	Subject   string `json:"sub"`
	Name      string `json:"name"`
	Username  string `json:"preferred_username"`
	Email     string `json:"email"`
	AvatarUrl string `json:"picture"`
}

type Error struct {
	ErrorCode        string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func getVerifierKey(state string) string {
	return "pkce_verifier:state:" + state
}

func getAccessTokenKey(sessionId string) string {
	return "access_token:session" + sessionId
}

type authService struct {
	store interfaces.Store
	keycloakURL, keycloakRealm,
	keycloakClientId, keycloakClientSecret,
	keycloakRedirectURI string
	encryptionKey string
}

func NewAuthService(store interfaces.Store,
	keycloakURL, keycloakRealm string,
	keycloakClientId, keycloakClientSecret string,
	keycloakRedirectURI, encryptionKey string) interfaces.AuthService {
	return &authService{
		store:                store,
		keycloakURL:          keycloakURL,
		keycloakRealm:        keycloakRealm,
		keycloakClientId:     keycloakClientId,
		keycloakClientSecret: keycloakClientSecret,
		keycloakRedirectURI:  keycloakRedirectURI,
		encryptionKey:        encryptionKey,
	}
}

func (s *authService) RedirectLogin(ctx context.Context) (string, error) {

	verifier, err := utils.CreateCodeVerifier()
	if err != nil {
		return "", fmt.Errorf("create code verifier: %w", err)
	}

	state := uuid.NewString()
	verifierKey := getVerifierKey(state)
	err = s.store.InMemory().Set(ctx, verifierKey, verifier.Value, time.Second*10)
	if err != nil {
		return "", fmt.Errorf("store in-memory set: %w", err)
	}

	challenge := verifier.CodeChallengeSHA256()
	redirectUrl := fmt.Sprintf(
		"%s/realms/%s/protocol/openid-connect/auth?"+
			"scope=openid"+
			"&response_type=code"+
			"&client_id=%s"+
			"&redirect_uri=%s"+
			"&state=%s"+
			"&code_challenge_method=S256"+
			"&code_challenge=%s",
		s.keycloakURL, s.keycloakRealm,
		s.keycloakClientId, s.keycloakRedirectURI,
		state, challenge)

	return redirectUrl, nil
}

func (s *authService) getToken(payload url.Values) (*Token, error) {
	tokenEndpoint := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", s.keycloakURL, s.keycloakRealm)
	res, err := http.PostForm(tokenEndpoint, payload)
	if err != nil {
		return nil, fmt.Errorf("http post form: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var KCErrorResponse Error
		if err := json.NewDecoder(res.Body).Decode(&KCErrorResponse); err != nil {
			return nil, fmt.Errorf("token request status: %d", res.StatusCode)
		}

		return nil, fmt.Errorf("token request status not OK: %s", KCErrorResponse.ErrorDescription)
	}

	var tokenResponse Token
	if err := json.NewDecoder(res.Body).Decode(&tokenResponse); err != nil {
		return nil, fmt.Errorf("token response decode: %w", err)
	}

	return &tokenResponse, nil
}

func (s *authService) getUserInfo(accessToken string) (*UserInfo, error) {
	userinfoEndpoint := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/userinfo", s.keycloakURL, s.keycloakRealm)
	req, err := http.NewRequest("GET", userinfoEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("http new request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http default client do: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var KCErrorResponse Error
		if err := json.NewDecoder(res.Body).Decode(&KCErrorResponse); err != nil {
			return nil, fmt.Errorf("user info request status: %d", res.StatusCode)
		}

		return nil, fmt.Errorf("user info request: %s", KCErrorResponse.ErrorDescription)
	}

	var keycloakUserInfo UserInfo
	if err := json.NewDecoder(res.Body).Decode(&keycloakUserInfo); err != nil {
		return nil, fmt.Errorf("user info response decode: %w", err)
	}

	return &keycloakUserInfo, nil
}

func (s *authService) Callback(ctx context.Context,
	state, code string,
	userAgent, device, ipAddress string) (*domain.Session, error) {

	verifierKey := utils.GetVerifierKey(state)
	state, err := s.store.InMemory().Get(ctx, verifierKey)
	if err != nil {
		return nil, fmt.Errorf("store in-memoty get: %w", err)
	}

	token, err := s.getToken(url.Values{
		"grant_type":    []string{"authorization_code"},
		"code":          []string{code},
		"code_verifier": []string{state},
		"redirect_uri":  []string{s.keycloakRedirectURI},
		"client_id":     []string{s.keycloakClientId},
		"client_secret": []string{s.keycloakClientSecret},
	})
	if err != nil {
		return nil, fmt.Errorf("get token: %w", err)
	}

	userInfo, err := s.getUserInfo(token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("get user info: %w", err)
	}

	var user *domain.User
	user, err = s.store.User().GetBySubject(ctx, userInfo.Subject, DEFAULT_PROVIDER)
	if err != nil {
		switch err {
		case apierr.ErrNotFound:
			userId, err := s.store.User().Create(ctx,
				userInfo.Subject, DEFAULT_PROVIDER,
				userInfo.Username, userInfo.Email, &userInfo.Name,
				&userInfo.AvatarUrl)
			if err != nil {
				return nil, fmt.Errorf("store user create: %w", err)
			}

			user, err = s.store.User().Get(ctx, userId)
			if err != nil {
				return nil, fmt.Errorf("store user get: %w", err)
			}
		}
	}

	expiresAt := time.Now().Add(time.Duration(token.RefreshExpiresIn) * time.Second)
	encryptedToken, err := utils.EncryptAES([]byte(token.RefreshToken), []byte(s.encryptionKey))
	if err != nil {
		return nil, fmt.Errorf("utils encrypt aes: %w", err)
	}

	sessionId, err := s.store.Session().Create(ctx,
		user.Id, encryptedToken,
		userAgent, ipAddress, device,
		expiresAt)
	if err != nil {
		return nil, fmt.Errorf("store session create: %w", err)
	}

	tokenKey := getAccessTokenKey(sessionId)
	err = s.store.InMemory().Set(ctx, tokenKey, token.AccessToken, time.Until(expiresAt))
	if err != nil {
		return nil, fmt.Errorf("store in-memoty set: %w", err)
	}

	session, err := s.store.Session().Get(ctx, sessionId)
	if err != nil {
		return nil, fmt.Errorf("store session get: %w", err)
	}

	return session, nil
}

func (s *authService) Authenticate(ctx context.Context,
	sessionId string) (string, error) {

	tokenKey := utils.GetAccessTokenKey(sessionId)
	accessToken, err := s.store.InMemory().Get(ctx, tokenKey)
	if err != nil {
		return "", fmt.Errorf("store in-memory get: %w", err)
	}

	sub, err := s.verifyAccessToken(ctx, accessToken)
	if err != nil {
		return "", fmt.Errorf("verify access token: %w", err)
	}

	user, err := s.store.User().GetBySubject(ctx, sub, DEFAULT_PROVIDER)
	if err != nil {
		return "", fmt.Errorf("store user get by subject: %w", err)
	}

	return user.Id, nil
}

func (s *authService) verifyAccessToken(ctx context.Context,
	token string) (string, error) {

	jwksKeySet, err := jwk.Fetch(ctx,
		fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs",
			s.keycloakURL, s.keycloakRealm))
	if err != nil {
		return "", fmt.Errorf("jwk fetch: %w", err)
	}

	parsedToken, err := jwt.Parse([]byte(token),
		jwt.WithKeySet(jwksKeySet),
		jwt.WithValidate(true))
	if err != nil {
		return "", fmt.Errorf("jwt parse: %w", err)
	}

	if parsedToken == nil {
		return "", fmt.Errorf("parsed token is null")
	}

	return parsedToken.Subject(), nil
}

func (s *authService) Refresh(ctx context.Context,
	sessionId string) error {

	session, err := s.store.Session().Get(ctx, sessionId)
	if err != nil {
		return fmt.Errorf("store session get: %w", err)
	}

	refreshToken, err := utils.DecryptAES(session.RefreshTokenEncrypted, []byte(s.encryptionKey))
	if err != nil {
		return fmt.Errorf("decrypt aes: %w", err)
	}

	tokenEndpoint := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", s.keycloakURL, s.keycloakRealm)
	res, err := http.PostForm(tokenEndpoint, url.Values{
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{string(refreshToken)},
		"redirect_uri":  []string{s.keycloakRedirectURI},
		"client_id":     []string{s.keycloakClientId},
		"client_secret": []string{s.keycloakClientSecret},
		"scope":         []string{"openid"},
	})
	if err != nil {
		return fmt.Errorf("http post form: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		// revoke session and remove from cookie

		err := s.store.Session().Revoke(ctx, sessionId)
		if err != nil {
			return fmt.Errorf("session revoke: %w", err)
		}

		var KCErrorResponse Error
		if err := json.NewDecoder(res.Body).Decode(&KCErrorResponse); err != nil {
			return fmt.Errorf("refresh token request decode: %w", err)
		}

		return fmt.Errorf("refresh token request: %w", err)
	}

	var token Token
	if err := json.NewDecoder(res.Body).Decode(&token); err != nil {
		return fmt.Errorf("new token decode: %w", err)
	}

	tokenExpiresAt := time.Now().Add(time.Duration(token.RefreshExpiresIn * int(time.Second)))
	encryptedRefreshToken, err := utils.EncryptAES([]byte(token.RefreshToken), []byte(s.encryptionKey))
	if err != nil {
		return fmt.Errorf("encrypt aes: %w", err)
	}

	err = s.store.Session().Update(ctx,
		sessionId,
		encryptedRefreshToken,
		tokenExpiresAt)
	if err != nil {
		return fmt.Errorf("session update: %w", err)
	}

	tokenKey := getAccessTokenKey(sessionId)
	err = s.store.InMemory().Set(ctx,
		tokenKey,
		token.AccessToken,
		time.Until(tokenExpiresAt))
	if err != nil {
		return fmt.Errorf("store in-memoty set: %w", err)
	}

	return nil
}

func (s *authService) Logout(ctx context.Context,
	sessionId string) error {
	err := s.store.Session().Revoke(ctx, sessionId)
	if err != nil {
		return fmt.Errorf("store session revoke: %w", err)
	}

	key := getAccessTokenKey(sessionId)
	err = s.store.InMemory().Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("store in-memory delete: %w", err)
	}

	return nil
}
