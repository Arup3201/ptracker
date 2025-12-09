package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	// migration
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		PG_HOST, PG_PORT, PG_USER, PG_PASS, PG_DB))
	if err != nil {
		fmt.Printf("[ERROR] postgres Open: %s", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		fmt.Printf("[ERROR] postgres WithInstance: %s", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		fmt.Printf("[ERROR] postgres NewWithDatabaseInstance: %s", err)
	}
	if err = m.Up(); err != nil {
		fmt.Printf("[ERROR] postgres migration error: %s", err)
	}

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

	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("[ERROR] server failed to start: %s", err)
	}
}
