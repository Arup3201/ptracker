package openid

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ptracker/core"
	"github.com/ptracker/core/users"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

const (
	USERINFO_ENDPOINT = "https://www.googleapis.com/oauth2/v3/userinfo"
)

// https://docs.cloud.google.com/identity-platform/docs/reference/rest/v1/UserInfo
type GoogleUserInfo struct {
	Subject string `json:"sub"`
	Name    string `json:"name"`
	Email   string `json:"email"`
}

type GoogleService struct {
	config    *oauth2.Config
	txManager *core.TxManager
	userRepo  *users.UserRepository
	oauthRepo *OauthRepository
}

func NewGoogleService(
	clientID, clientSecret string,
	redirectURI string,
	txManager *core.TxManager,
	userRepo *users.UserRepository,
	oauthRepo *OauthRepository,
) *GoogleService {

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"openid", "email"},
		Endpoint:     google.Endpoint,
		RedirectURL:  redirectURI,
	}

	return &GoogleService{
		config:    conf,
		txManager: txManager,
		userRepo:  userRepo,
		oauthRepo: oauthRepo,
	}
}

func (s *GoogleService) GetAuthCodeURL(ctx context.Context) string {
	// TODO: generate state and verifier and store it in redis
	// verifier := oauth2.GenerateVerifier()
	url := s.config.AuthCodeURL("secret-state")
	return url
}

func (s *GoogleService) GetUserInfoFromAuthCode(ctx context.Context,
	state, code string) (*GoogleUserInfo, error) {

	if state != "secret-state" {
		return nil, fmt.Errorf("state mismatch: %w", core.ErrInvalidValue)
	}

	tok, err := s.config.Exchange(ctx, code)
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

func (s *GoogleService) GetUserID(ctx context.Context,
	subject, provider string) (string, error) {

	acc, err := s.oauthRepo.Get(ctx, subject, provider)
	if err != nil {
		return "", fmt.Errorf("oauth repository get: %w", err)
	}

	return acc.UserID, nil
}

func (s *GoogleService) CreateAccount(ctx context.Context,
	userInfo GoogleUserInfo) (string, error) {

	var err error
	var userID string

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
