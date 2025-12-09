package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ptracker/handlers"
)

const (
	HOST = "localhost"
	PORT = 8081
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/keycloak/login", handlers.KeycloakLogin)
	mux.HandleFunc("GET /api/keycloak/callback", handlers.KeycloakCallback)

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
