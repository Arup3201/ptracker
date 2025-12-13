package main

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ptracker/db"
	"github.com/ptracker/handlers"
)

const (
	HOST = "localhost"
	PORT = 8081

	PG_HOST = "127.0.0.1"
	PG_PORT = "5432"
	PG_USER = "postgres"
	PG_PASS = "1234"
	PG_DB   = "ptracker"
)

func main() {
	// DB connection
	db.ConnectPostgres(fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", PG_HOST, PG_PORT,
		PG_USER, PG_PASS, PG_DB))

	// migrate
	db.Migrate()

	// handler
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/keycloak/login", handlers.KeycloakLogin)
	mux.HandleFunc("GET /api/keycloak/callback", handlers.KeycloakCallback)

	// server
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
