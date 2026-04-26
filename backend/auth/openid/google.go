package openid

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ptracker/auth"
	"github.com/ptracker/core"
	"github.com/ptracker/core/users"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

const (
	USERINFO_ENDPOINT = "https://www.googleapis.com/oauth2/v3/userinfo"
)

func getStateKey(s string) string {
	return "AUTH_STATE:" + s
}

func getVerifierKey(s string) string {
	return "PKCE_VERIFIER:" + s
}

// https://docs.cloud.google.com/identity-platform/docs/reference/rest/v1/UserInfo
type GoogleUserInfo struct {
	Subject string `json:"sub"`
	Name    string `json:"name"`
	Email   string `json:"email"`
}

type GoogleService struct {
	config      *oauth2.Config
	txManager   *core.TxManager
	userRepo    *users.UserRepository
	oauthRepo   *OauthRepository
	stringStore *StringStore
}

func NewGoogleService(
	clientID, clientSecret string,
	redirectURI string,
	txManager *core.TxManager,
	userRepo *users.UserRepository,
	oauthRepo *OauthRepository,
	stringStore *StringStore,
) *GoogleService {

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"openid", "email"},
		Endpoint:     google.Endpoint,
		RedirectURL:  redirectURI,
	}

	return &GoogleService{
		config:      conf,
		txManager:   txManager,
		userRepo:    userRepo,
		oauthRepo:   oauthRepo,
		stringStore: stringStore,
	}
}

func (s *GoogleService) GetAuthCodeURL(ctx context.Context) (string, error) {

	state, _ := auth.GetRandomToken(32)
	err := s.stringStore.Store(ctx, getStateKey(state), state, 1*time.Minute)
	if err != nil {
		return "", fmt.Errorf("string store Store: %w", err)
	}
	verifier := oauth2.GenerateVerifier()
	err = s.stringStore.Store(ctx, getVerifierKey(state), verifier, 1*time.Minute)
	if err != nil {
		return "", fmt.Errorf("string store Store: %w", err)
	}
	url := s.config.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier))
	return url, nil
}

func (s *GoogleService) getUserInfo(ctx context.Context,
	state, code string) (*GoogleUserInfo, error) {

	_, err := s.stringStore.Get(ctx, getStateKey(state))
	if err != nil {
		return nil, fmt.Errorf("stringStore Get state: %w", err)
	}

	verifier, err := s.stringStore.Get(ctx, getVerifierKey(state))
	if err != nil {
		return nil, fmt.Errorf("stringStore Get PKCE verifier: %w", err)
	}

	tok, err := s.config.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		return nil, fmt.Errorf("code exchange: %w", err)
	}

	client := s.config.Client(ctx, tok)
	res, err := client.Get(USERINFO_ENDPOINT)
	if err != nil {
		return nil, fmt.Errorf("client get: %w", err)
	}
	defer res.Body.Close()

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("decode userinfo: %w", err)
	}

	return &userInfo, nil
}

func (s *GoogleService) GetUser(ctx context.Context,
	subject, provider string) (string, error) {

	acc, err := s.oauthRepo.Get(ctx, subject, provider)
	if err != nil {
		return "", fmt.Errorf("oauth repository get: %w", err)
	}

	return acc.UserID, nil
}

func (s *GoogleService) CreateAccount(ctx context.Context,
	state, code string) (string, error) {

	var err error
	var userID string

	userInfo, err := s.getUserInfo(ctx, state, code)
	if err != nil {
		return "", fmt.Errorf("getUserInfo: %w", err)
	}

	var username = strings.Split(userInfo.Email, "@")[0]

	err = s.txManager.WithTx(func(tx *gorm.DB) error {
		userRepo := s.userRepo.WithTx(tx)
		oauthRepo := s.oauthRepo.WithTx(tx)

		userID, err = userRepo.Create(
			ctx,
			username,
			userInfo.Email,
			&userInfo.Name,
			nil,
		)
		if err != nil {
			return fmt.Errorf("user repository create: %w", err)
		}

		err = oauthRepo.Create(
			ctx,
			userInfo.Subject,
			OAUTH_PROVIDER_GOOGLE,
			userID,
			userInfo.Email,
		)
		if err != nil {
			return fmt.Errorf("oauth account repository create: %w", err)
		}

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("transaction: %w", err)
	}

	return userID, nil
}
