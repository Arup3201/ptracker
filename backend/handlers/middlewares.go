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

var accessTokens = map[string]string{}

func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		public := []string{"/auth/login", "/auth/callback", "/auth/refresh"}
		for _, endpoint := range public {
			if strings.Contains(strings.TrimPrefix(r.URL.Path, "/api"), endpoint) {
				next.ServeHTTP(w, r)
				return
			}
		}

		cookies := r.Cookies()
		ind := slices.IndexFunc(cookies, func(cookie *http.Cookie) bool {
			return cookie.Name == SESSION_COOKIE_NAME
		})
		if ind == -1 {
			fmt.Printf("[ERROR] authorization error: session cookie missing\n")

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ApiResponse{
				Status:  "Error",
				Message: "Unathorized",
			})
			return
		}

		sessionId := cookies[ind].Value
		sub, err := verifyAccessToken(accessTokens[sessionId])
		if err != nil {

			fmt.Printf("[ERROR] authorization error: refresh token error: %s\n", err)

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ApiResponse{
				Status:  "Error",
				Message: "Unathorized",
			})
			return
		}

		fmt.Printf("subject: %s\n", sub)

		next.ServeHTTP(w, r)
	})
}

func verifyAccessToken(accessToken string) (string, error) {
	jwksKeySet, err := jwk.Fetch(context.TODO(), fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", KC_URL, KC_REALM))
	if err != nil {
		return "", err
	}

	token, err := jwt.Parse([]byte(accessToken), jwt.WithKeySet(jwksKeySet), jwt.WithValidate(true))
	if err != nil {
		return "", err
	}

	if token == nil {
		return "", fmt.Errorf("parsed token is null")
	}

	return token.Subject(), nil
}
