package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ptracker/internal/constants"
	"github.com/ptracker/internal/interfaces"
	"github.com/ptracker/internal/utils"
)

type authController struct {
	authService interfaces.AuthService
	homeURL     string
}

func NewAuthController(service interfaces.AuthService,
	homeURL string) *authController {
	return &authController{
		authService: service,
		homeURL:     homeURL,
	}
}

func (c *authController) Login(w http.ResponseWriter, r *http.Request) error {

	url, err := c.authService.RedirectLogin(r.Context())
	if err != nil {
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Login failed",
			Err:     fmt.Errorf("keycloak login: redis set state: %w", err),
		}
	}

	http.Redirect(w, r, url, http.StatusSeeOther)

	return nil
}

func (c *authController) Callback(w http.ResponseWriter, r *http.Request) error {

	authorization_code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	authorization_error_code := r.URL.Query().Get("error")

	if authorization_error_code != "" {
		authorization_error_description := r.URL.Query().Get("error_description")
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Authorization denied",
			Err:     fmt.Errorf("authorization code error: %s", authorization_error_description),
		}
	}

	if state == "" {
		return &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Authorization denied",
			Err:     fmt.Errorf("state missing"),
		}
	}

	userAgent := r.UserAgent()
	device := utils.ParseUserAgent(userAgent)
	ipAddress := strings.Split(r.RemoteAddr, ":")[0]

	session, err := c.authService.Callback(r.Context(),
		state, authorization_code,
		userAgent, device, ipAddress)
	if err != nil {
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Callback failed",
			Err:     fmt.Errorf("auth service callback: %w", err),
		}
	}

	cookie := &http.Cookie{
		Name:     constants.SESSION_COOKIE_NAME,
		Value:    session.Id,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		Expires:  session.ExpiresAt,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, c.homeURL, http.StatusSeeOther)

	return nil
}

func (c *authController) Refresh(w http.ResponseWriter, r *http.Request) error {

	cookies := r.Cookies()
	sessionId, err := utils.GetSessionIdFromCookie(cookies, constants.SESSION_COOKIE_NAME)
	if err != nil {
		return &HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "User session has expired",
			Err:     fmt.Errorf("utils get session id from cookie: %w", err),
		}
	}

	err = c.authService.Refresh(r.Context(), sessionId)

	if err != nil {
		// remove session from cookie

		cookie := &http.Cookie{
			Name:     constants.SESSION_COOKIE_NAME,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteDefaultMode,
			Expires:  time.Unix(0, 0),
		}
		http.SetCookie(w, cookie)

		return &HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "Session cannot be refreshed. Try to login again.",
			Err:     fmt.Errorf("auth service refresh: %w", err),
		}
	}

	json.NewEncoder(w).Encode(HTTPSuccessResponse[any]{
		Status:  RESPONSE_SUCCESS_STATUS,
		Message: "Access token refreshed",
	})
	return nil
}

func (c *authController) Logout(w http.ResponseWriter, r *http.Request) error {

	sessionId, err := utils.GetSessionIdFromCookie(r.Cookies(), constants.SESSION_COOKIE_NAME)
	if err != nil {
		return &HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "User session has expired",
			Err:     fmt.Errorf("utils get session id from cookie: %w", err),
		}
	}

	err = c.authService.Logout(r.Context(), sessionId)
	if err != nil {
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Logout error, please try later",
			Err:     fmt.Errorf("auth service logout: %w", err),
		}
	}

	cookie := &http.Cookie{
		Name:     constants.SESSION_COOKIE_NAME,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		Expires:  time.Unix(0, 0),
	}
	http.SetCookie(w, cookie)

	json.NewEncoder(w).Encode(HTTPSuccessResponse[any]{
		Status:  RESPONSE_SUCCESS_STATUS,
		Message: "Logout success",
	})
	return nil
}
