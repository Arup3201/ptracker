package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ptracker/auth"
	"github.com/ptracker/auth/openid"
	"github.com/ptracker/core"
	"github.com/ptracker/core/users"
)

type GoogleApi struct {
	googleService                     *openid.GoogleService
	tokenService                      *auth.TokenService
	userService                       *users.UserService
	frontendHomeURL, frontendLoginURL string
}

func NewGoogleApi(
	googleService *openid.GoogleService,
	tokenService *auth.TokenService,
	userService *users.UserService,
	homeURL, loginURL string,
) *GoogleApi {
	return &GoogleApi{
		googleService:    googleService,
		tokenService:     tokenService,
		userService:      userService,
		frontendHomeURL:  homeURL,
		frontendLoginURL: loginURL,
	}
}

func (api *GoogleApi) Redirect(w http.ResponseWriter, r *http.Request) error {

	url, err := api.googleService.GetAuthCodeURL(r.Context())
	if err != nil {
		return fmt.Errorf("google service GetAuthCodeURL: %w", err)
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[string]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data:   &url,
	})

	return nil
}

func (api *GoogleApi) Callback(w http.ResponseWriter, r *http.Request) error {

	errParam := r.URL.Query().Get("error")
	if errParam != "" {
		http.Redirect(w, r, api.frontendLoginURL+"?error="+errParam, http.StatusSeeOther)
		return nil
	}

	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	if state == "" {
		http.Redirect(w, r, api.frontendLoginURL+"?error=missing_state", http.StatusSeeOther)
		return nil
	}

	if code == "" {
		http.Redirect(w, r, api.frontendLoginURL+"?error=missing_auth_code", http.StatusSeeOther)
		return nil
	}

	userID, token, err := api.googleService.Callback(
		r.Context(),
		state,
		code,
	)
	switch err {
	case core.ErrInvalidValue:
		http.Redirect(w, r, api.frontendLoginURL+"?error=invalid_state_or_code", http.StatusSeeOther)
		return nil
	case core.ErrDuplicate:
		http.Redirect(w, r, api.frontendLoginURL+"?error=account_exist", http.StatusSeeOther)
		return nil
	}
	if err != nil {
		http.Redirect(w, r, api.frontendLoginURL+"?error=server_error", http.StatusSeeOther)
		return nil
	}

	http.Redirect(w, r,
		api.frontendLoginURL+"?user_id="+userID+"&token="+token,
		http.StatusSeeOther)

	return nil
}

func (api *GoogleApi) Login(w http.ResponseWriter, r *http.Request) error {

	userID := r.URL.Query().Get("user_id")
	token := r.URL.Query().Get("token")
	if userID == "" || token == "" {
		return core.ErrInvalidValue
	}

	err := api.googleService.ValidToken(
		r.Context(),
		token,
	)
	if err != nil {
		return err
	}

	refreshToken, err := api.tokenService.CreateRefreshToken(
		r.Context(),
		userID,
	)
	if err != nil {
		return fmt.Errorf("token service create refresh token: %w", err)
	}

	accessToken, err := api.tokenService.CreateAccessToken(
		r.Context(),
		userID,
	)
	if err != nil {
		return fmt.Errorf("token service create access token: %w", err)
	}

	user, err := api.userService.Get(
		r.Context(),
		userID)
	if err != nil {
		return fmt.Errorf("user service get: %w", err)
	}

	cookie := &http.Cookie{
		Name:     REFRESH_TOKEN_COOKIE_NAME,
		Value:    refreshToken.Value,
		Path:     "/", // TODO: auth path only
		Expires:  refreshToken.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)

	json.NewEncoder(w).Encode(HTTPSuccessResponse[LoginResponse]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data: &LoginResponse{
			AccessToken: accessToken.Value,
			ExpiresAt:   accessToken.ExpiresAt,
			User: core.Avatar{
				UserID:      user.ID,
				Username:    user.Username,
				Email:       user.Email,
				DisplayName: user.DisplayName,
				AvatarURL:   user.AvatarURL,
			},
		},
	})

	return nil
}
