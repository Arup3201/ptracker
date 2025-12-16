package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
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

type HTTPErrorHandler func(w http.ResponseWriter, r *http.Request) error

func (fn HTTPErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		fmt.Printf("[ERROR] %s\n", err)

		if httpError, ok := err.(*HTTPError); ok {
			w.WriteHeader(httpError.Code)
			json.NewEncoder(w).Encode(HTTPErrorResponse{
				Status:  "error",
				Message: httpError.Message,
			})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(HTTPErrorResponse{
				Status:  "error",
				Message: "Unexpected error occured, we are working on it.",
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

		cookies := r.Cookies()
		ind := slices.IndexFunc(cookies, func(cookie *http.Cookie) bool {
			return cookie.Name == SESSION_COOKIE_NAME
		})
		if ind == -1 {
			return &HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "User is not authorized",
				Err:     fmt.Errorf("authorize: session cookie not found"),
			}
		}

		sessionId := cookies[ind].Value
		sub, err := verifyAccessToken(accessTokens[sessionId])
		if err != nil {
			return &HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "User is not authorized",
				Err:     fmt.Errorf("authorize: %w", err),
			}
		}

		fmt.Printf("subject: %s\n", sub)

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
