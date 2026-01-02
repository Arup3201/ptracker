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
	"github.com/redis/go-redis/v9"
)

var (
	HOST              string
	PORT              string
	PG_HOST           string
	PG_PORT           string
	PG_USER           string
	PG_PASS           string
	PG_DB             string
	KC_URL            string
	KC_REALM          string
	KC_CLIENT_ID      string
	KC_CLIENT_SECRET  string
	KC_REDIRECT_URI   string
	ENCRYPTION_SECRET string
	HOME_URL          string
)

func getEnvironment() error {
	HOST = os.Getenv("HOST")
	if HOST == "" {
		return fmt.Errorf("environment variable 'HOST' missing")
	}
	PORT = os.Getenv("PORT")
	if PORT == "" {
		return fmt.Errorf("environment variable 'PORT' missing")
	}
	HOME_URL = os.Getenv("HOME_URL")
	if HOME_URL == "" {
		return fmt.Errorf("environment variable 'HOME_URL' missing")
	}
	ENCRYPTION_SECRET = os.Getenv("ENCRYPTION_SECRET")
	if ENCRYPTION_SECRET == "" {
		return fmt.Errorf("environment variable 'ENCRYPTION_SECRET' missing")
	}
	PG_HOST = os.Getenv("PG_HOST")
	if PG_HOST == "" {
		return fmt.Errorf("environment variable 'PG_HOST' missing")
	}
	PG_USER = os.Getenv("PG_USER")
	if PG_USER == "" {
		return fmt.Errorf("environment variable 'PG_USER' missing")
	}
	PG_PORT = os.Getenv("PG_PORT")
	if PG_PORT == "" {
		return fmt.Errorf("environment variable 'PG_PORT' missing")
	}
	PG_PASS = os.Getenv("PG_PASS")
	if PG_PASS == "" {
		return fmt.Errorf("environment variable 'PG_PASS' missing")
	}
	PG_DB = os.Getenv("PG_DB")
	if PG_DB == "" {
		return fmt.Errorf("environment variable 'PG_DB' missing")
	}
	KC_URL = os.Getenv("KC_URL")
	if KC_URL == "" {
		return fmt.Errorf("environment variable 'KC_URL' missing")
	}
	KC_REALM = os.Getenv("KC_REALM")
	if KC_REALM == "" {
		return fmt.Errorf("environment variable 'KC_REALM' missing")
	}
	KC_CLIENT_ID = os.Getenv("KC_CLIENT_ID")
	if KC_CLIENT_ID == "" {
		return fmt.Errorf("environment variable 'KC_CLIENT_ID' missing")
	}
	KC_CLIENT_SECRET = os.Getenv("KC_CLIENT_SECRET")
	if KC_CLIENT_SECRET == "" {
		return fmt.Errorf("environment variable 'KC_CLIENT_SECRET' missing")
	}
	KC_REDIRECT_URI = os.Getenv("KC_REDIRECT_URI")
	if KC_REDIRECT_URI == "" {
		return fmt.Errorf("environment variable 'KC_REDIRECT_URI' missing")
	}

	return nil
}

func attachMiddlewares(
	mux *http.ServeMux,
	redis *redis.Client,
	pattern string,
	handler handlers.HTTPErrorHandler) {
	authMiddleware := handlers.Authorize(redis, KC_URL, KC_REALM)
	mux.Handle(pattern, handlers.HTTPErrorHandler(authMiddleware(handler)))
}

func main() {
	// Get constants from environment
	err := getEnvironment()
	if err != nil {
		log.Fatalf("[ERROR] server failed to get environemnt: %s", err)
	}

	// DB connection
	err = db.ConnectPostgres(fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", PG_HOST, PG_PORT,
		PG_USER, PG_PASS, PG_DB))
	if err != nil {
		log.Fatalf("[ERROR] server failed to connect to postgres: %s", err)
	}

	// migrate
	err = db.Migrate()
	if err != nil {
		log.Fatalf("[ERROR] server failed to migrate postgres: %s", err)
	}

	// Redis
	redis := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	// handler
	mux := http.NewServeMux()

	kcHandler, err := handlers.CreateKeycloakHandler(
		KC_URL, KC_REALM, KC_CLIENT_ID, KC_CLIENT_SECRET, KC_REDIRECT_URI,
		HOME_URL, ENCRYPTION_SECRET, redis,
	)
	if err != nil {
		log.Fatalf("[ERROR] server failed to create keycloak handler: %s", err)
	}

	attachMiddlewares(mux, redis, "GET /api/auth/login", kcHandler.KeycloakLogin)
	attachMiddlewares(mux, redis, "GET /api/auth/callback", kcHandler.KeycloakCallback)
	attachMiddlewares(mux, redis, "POST /api/auth/refresh", kcHandler.KeycloakRefresh)
	attachMiddlewares(mux, redis, "POST /api/auth/logout", kcHandler.KeycloakLogout)

	rateLimiter := handlers.TokenBucketRateLimiter(redis, 5, 2)
	attachMiddlewares(mux, redis, "POST /api/projects", rateLimiter(handlers.CreateProject))

	// server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", HOST, PORT),
		Handler:      handlers.Logging(mux),
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("[INFO] server starting at %s:%s\n", HOST, PORT)

	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("[ERROR] server failed to start: %s", err)
	}
}
