package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/ptracker/auth"
	"github.com/ptracker/auth/openid"
	"github.com/ptracker/core"
	"github.com/ptracker/core/users"
)

type GoogleApi struct {
	googleService *openid.GoogleService
	tokenService  *auth.TokenService
	userService   *users.UserService
}

func NewGoogleApi(
	googleService *openid.GoogleService,
	tokenService *auth.TokenService,
	userService *users.UserService,
) *GoogleApi {
	return &GoogleApi{
		googleService: googleService,
		tokenService:  tokenService,
		userService:   userService,
	}
}

func (api *GoogleApi) Redirect(w http.ResponseWriter, r *http.Request) error {

	url := api.googleService.GetAuthCodeURL(r.Context())

	json.NewEncoder(w).Encode(HTTPSuccessResponse[string]{
		Status: RESPONSE_SUCCESS_STATUS,
		Data:   &url,
	})

	return nil
}

func (api *GoogleApi) Callback(w http.ResponseWriter, r *http.Request) error {

	errParam := r.URL.Query().Get("error")

	var responseScript, userID string
	var userInfo *openid.GoogleUserInfo
	var accessToken, refreshToken *auth.Token
	var user *users.User
	var err error
	if errParam == "" {
		code := r.FormValue("code")
		state := r.FormValue("state")

		userInfo, err = api.googleService.GetUserInfoFromAuthCode(
			r.Context(),
			state,
			code,
		)
		if err == nil {
			userID, err = api.googleService.GetUserID(
				r.Context(),
				userInfo.Subject,
				openid.OAUTH_PROVIDER_GOOGLE,
			)
			if errors.Is(err, core.ErrNotFound) {
				userID, err = api.googleService.CreateAccount(
					r.Context(),
					*userInfo,
				)
				if err != nil {
					log.Printf("[ERROR] google service create account: %s", err)

					responseScript = `
					window.opener.postMessage(
						{ success: false, error: "Email may not be verified or user already exist with another method" }
					);
					`
				}
			}
			if err == nil {
				// nil err for GetUserID and CreateAccount

				refreshToken, err = api.tokenService.CreateRefreshToken(
					r.Context(),
					userID,
				)
				if err == nil {
					accessToken, err = api.tokenService.CreateAccessToken(
						r.Context(),
						userID,
					)
					if err == nil {
						user, err = api.userService.Get(
							r.Context(),
							userID)
						if err == nil {
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

							data, _ := json.Marshal(core.Avatar{
								UserID:      user.ID,
								Username:    user.Username,
								Email:       user.Email,
								DisplayName: user.DisplayName,
								AvatarURL:   user.AvatarURL,
							})

							responseScript = `
								window.opener.postMessage(
									{ success: true,
									access_token: "` + accessToken.Value + `",
									expires_at: "` + accessToken.ExpiresAt.Format(time.RFC3339) + `",
									user: "` + string(data) + `"
									}
								);
							`
						} else {
							log.Printf("[ERROR] user service get: %s", err)

							responseScript = `
							window.opener.postMessage(
								{ success: false, error: "Failed to get user data" }
							);
							`
						}
					} else {
						log.Printf("[ERROR] token service create access token: %s", err)

						responseScript = `
							window.opener.postMessage(
								{ success: false, error: "Failed to create access token" }
							);
							`
					}
				} else {
					log.Printf("[ERROR] token service create refresh token: %s", err)

					responseScript = `
							window.opener.postMessage(
								{ success: false, error: "Failed to create refresh token" }
							);
							`
				}
			} else if !errors.Is(err, core.ErrNotFound) {
				log.Printf("[ERROR] google service get user ID: %s", err)

				// GetUserID failed for some reason
				responseScript = `
					window.opener.postMessage(
						{ success: false, error: "Server error" }
					);
					`
			}
		} else {
			log.Printf("[ERROR] google service get user info from AuthCode: %s", err)
			responseScript = `
					window.opener.postMessage(
						{ success: false, error: "Failed to get user profile information" }
					);
					`
		}
	} else {
		responseScript = `
		window.opener.postMessage(
			{ success: false, error: "` + errParam + `" }
		);
		`
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
		<html>
            <body>
                <script>
                    ` + responseScript + `
					window.close();
                </script>
            </body>
        </html>
	`))

	return nil
}
