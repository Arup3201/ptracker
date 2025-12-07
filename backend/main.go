package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	HOST            = "localhost"
	PORT            = 8081
	KC_URL          = "http://localhost:8080"
	KC_REALM        = "ptracker"
	KC_CLIENT_ID    = "api"
	KC_REDIRECT_URI = "http://localhost:8081/api/keycloak/callback"
)

type ApiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ApiData map[string]any

type ApiResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Error   ApiError `json:"error,omitempty"`
	Data    ApiData  `json:"data,omitempty"`
}

var verifier []byte

func KeycloakLogin(w http.ResponseWriter, r *http.Request) {
	rand.Read(verifier)

	h := sha256.New()
	h.Write(verifier)
	hashedVerifier := h.Sum(nil)

	kc_auth_url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/auth?"+
		"scope=openid"+
		"&response_type=code"+
		"&client_id=%s"+
		"&redirect_uri=%s"+
		"&code_challenge_method=S256"+
		"&code_challenge=%s",
		KC_URL, KC_REALM, KC_CLIENT_ID, KC_REDIRECT_URI, fmt.Sprintf("%x", hashedVerifier))
	http.Redirect(w, r, kc_auth_url, http.StatusSeeOther)
}

func KeycloakCallback(w http.ResponseWriter, r *http.Request) {
	authorization_code := r.URL.Query().Get("code")
	authorization_error_code := r.URL.Query().Get("error")

	if authorization_error_code != "" {
		authorization_error_description := r.URL.Query().Get("error_description")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ApiResponse{
			Status:  "Error",
			Message: "Keycloak authorization error",
			Error: ApiError{
				Code:    authorization_error_code,
				Message: authorization_error_description,
			},
		})
		return
	}

	json.NewEncoder(w).Encode(ApiResponse{
		Status:  "Success",
		Message: "Authorization success",
		Data: map[string]any{
			"code": authorization_code,
		},
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/keycloak/login", KeycloakLogin)
	mux.HandleFunc("GET /api/keycloak/callback", KeycloakCallback)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", HOST, PORT),
		Handler:      mux,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("[INFO] server starting at %s:%d\n", HOST, PORT)

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("[ERROR] server failed to start: %s", err)
	}
}
