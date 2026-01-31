package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ptracker/constants"
	"github.com/ptracker/controllers"
	"github.com/ptracker/interfaces"
	"github.com/ptracker/utils"
)

type authMiddleware struct {
	service interfaces.AuthService
}

func NewAuthMiddleware(service interfaces.AuthService) Middleware {
	return &authMiddleware{
		service: service,
	}
}

func (m *authMiddleware) Next(next http.Handler) controllers.HTTPErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		public := []string{"/auth/login", "/auth/callback", "/auth/refresh"}
		for _, endpoint := range public {
			if strings.Contains(strings.TrimPrefix(r.URL.Path, "/api"), endpoint) {
				next.ServeHTTP(w, r)
				return nil
			}
		}

		sessionId, err := utils.GetSessionIdFromCookie(r.Cookies(), constants.SESSION_COOKIE_NAME)
		if err != nil {
			return &controllers.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "User is not authorized",
				Err:     fmt.Errorf("authorize: session cookie not found"),
			}
		}

		userId, err := m.service.Authenticate(r.Context(), sessionId)
		if err != nil {
			return &controllers.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Server failed to authenticate",
				Err:     fmt.Errorf("auth service authenticate: %w", err),
			}
		}

		ctx := context.WithValue(r.Context(), "user_id", userId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
		return nil
	}
}
