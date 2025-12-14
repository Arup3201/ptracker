package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ptracker/db"
	"github.com/ptracker/handlers"
)

var (
	HOST string
	PORT string
)

func setConstants() {
	HOST = os.Getenv("HOST")
	if HOST == "" {
		log.Fatalf("environment variable 'HOST' missing")
	}
	PORT = os.Getenv("PORT")
	if PORT == "" {
		log.Fatalf("environment variable 'PORT' missing")
	}
	db.PG_USER = os.Getenv("PG_USER")
	if db.PG_USER == "" {
		log.Fatalf("environment variable 'PG_USER' missing")
	}
	db.PG_PORT = os.Getenv("PG_PORT")
	if db.PG_PORT == "" {
		log.Fatalf("environment variable 'PG_PORT' missing")
	}
	db.PG_PASS = os.Getenv("PG_PASS")
	if db.PG_PASS == "" {
		log.Fatalf("environment variable 'PG_PASS' missing")
	}
	db.PG_DB = os.Getenv("PG_DB")
	if db.PG_DB == "" {
		log.Fatalf("environment variable 'PG_DB' missing")
	}
	handlers.KC_URL = os.Getenv("KC_URL")
	if handlers.KC_URL == "" {
		log.Fatalf("environment variable 'KC_URL' missing")
	}
	handlers.KC_REALM = os.Getenv("KC_REALM")
	if handlers.KC_REALM == "" {
		log.Fatalf("environment variable 'KC_REALM' missing")
	}
	handlers.KC_CLIENT_ID = os.Getenv("KC_CLIENT_ID")
	if handlers.KC_CLIENT_ID == "" {
		log.Fatalf("environment variable 'KC_CLIENT_ID' missing")
	}
	handlers.KC_CLIENT_SECRET = os.Getenv("KC_CLIENT_SECRET")
	if handlers.KC_CLIENT_SECRET == "" {
		log.Fatalf("environment variable 'KC_CLIENT_SECRET' missing")
	}
	handlers.KC_REDIRECT_URI = os.Getenv("KC_REDIRECT_URI")
	if handlers.KC_REDIRECT_URI == "" {
		log.Fatalf("environment variable 'KC_REDIRECT_URI' missing")
	}
}

func main() {
	// Set constants from environment
	setConstants()

	// DB connection
	db.ConnectPostgres()

	// migrate
	db.Migrate()

	// handler
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/keycloak/login", handlers.KeycloakLogin)
	mux.HandleFunc("GET /api/keycloak/callback", handlers.KeycloakCallback)

	// server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", HOST, PORT),
		Handler:      mux,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("[INFO] server starting at %s:%s\n", HOST, PORT)

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("[ERROR] server failed to start: %s", err)
	}
}
