package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/ptracker/db"
	"github.com/ptracker/utils"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)

		log.Printf("%s %s - %d", r.Method, r.RequestURI, rec.Result().StatusCode)

		maps.Copy(w.Header(), rec.Header())
		w.WriteHeader(rec.Result().StatusCode)
		w.Write(rec.Body.Bytes())
	})
}

// HTTP Error ID
const (
	ERR_UNAUTHORIZED       = "unauthorized"
	ERR_INVALID_BODY       = "invalid_body"
	ERR_ACCESS_DENIED      = "access_denied"
	ERR_RESOURCE_NOT_FOUND = "resource_not_found"
	ERR_SERVER_ERROR       = "server_error"
)

type HTTPErrorHandler func(w http.ResponseWriter, r *http.Request) error

func (fn HTTPErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		fmt.Printf("[ERROR] %s\n", err)

		if httpError, ok := err.(*HTTPError); ok {
			w.WriteHeader(httpError.Code)
			switch httpError.Code {
			case http.StatusBadRequest:
				json.NewEncoder(w).Encode(HTTPErrorResponse{
					Status: "error",
					Error: ErrorBody{
						Id:      ERR_INVALID_BODY,
						Message: httpError.Message,
					},
				})
			case http.StatusUnauthorized:
				json.NewEncoder(w).Encode(HTTPErrorResponse{
					Status: "error",
					Error: ErrorBody{
						Id:      ERR_UNAUTHORIZED,
						Message: httpError.Message,
					},
				})
			case http.StatusForbidden:
				json.NewEncoder(w).Encode(HTTPErrorResponse{
					Status: "error",
					Error: ErrorBody{
						Id:      ERR_ACCESS_DENIED,
						Message: httpError.Message,
					},
				})
			case http.StatusNotFound:
				json.NewEncoder(w).Encode(HTTPErrorResponse{
					Status: "error",
					Error: ErrorBody{
						Id:      ERR_RESOURCE_NOT_FOUND,
						Message: httpError.Message,
					},
				})
			default:
				json.NewEncoder(w).Encode(HTTPErrorResponse{
					Status: "error",
					Error: ErrorBody{
						Id:      ERR_SERVER_ERROR,
						Message: httpError.Message,
					},
				})
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(HTTPErrorResponse{
				Status: "error",
				Error: ErrorBody{
					Id:      ERR_SERVER_ERROR,
					Message: "Something unexpected happened, please try again later.",
				},
			})
		}
		return
	}
}

var accessTokens = map[string]string{}

func Authorize(next http.Handler) HTTPErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		public := []string{"/auth/login", "/auth/callback", "/auth/refresh"}
		for _, endpoint := range public {
			if strings.Contains(strings.TrimPrefix(r.URL.Path, "/api"), endpoint) {
				next.ServeHTTP(w, r)
				return nil
			}
		}

		sessionId, err := utils.GetSessionIdFromCookie(r.Cookies(), SESSION_COOKIE_NAME)
		if err != nil {
			return &HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "User is not authorized",
				Err:     fmt.Errorf("authorize: session cookie not found"),
			}
		}

		sub, err := verifyAccessToken(accessTokens[sessionId])
		if err != nil {
			return &HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "User is not authorized",
				Err:     fmt.Errorf("authorize: %w", err),
			}
		}

		user, err := db.GetUserBySub(sub)

		ctx := context.WithValue(r.Context(), "user_id", user.Id)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
		return nil
	}
}

func verifyAccessToken(accessToken string) (string, error) {
	jwksKeySet, err := jwk.Fetch(context.TODO(), fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", KC_URL, KC_REALM))
	if err != nil {
		return "", fmt.Errorf("verify access token: %w", err)
	}

	token, err := jwt.Parse([]byte(accessToken), jwt.WithKeySet(jwksKeySet), jwt.WithValidate(true))
	if err != nil {
		return "", fmt.Errorf("verify access token: %w", err)
	}

	if token == nil {
		return "", fmt.Errorf("parsed token is null")
	}

	return token.Subject(), nil
}
