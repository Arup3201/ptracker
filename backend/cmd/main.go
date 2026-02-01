package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ptracker/internal"
	"github.com/ptracker/internal/infra"
	"github.com/redis/go-redis/v9"
)

func main() {

	config := &internal.Config{}
	config.Load()

	// DB connection
	database, err := infra.NewDatabase("postgres",
		fmt.Sprintf("host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=disable", config.DbHost, config.DbPort,
			config.DbUser, config.DbPass, config.DbName))
	if err != nil {
		log.Fatalf("[ERROR] server failed to connect to postgres: %s", err)
	}

	// Redis
	redis := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	inMemory := infra.NewInMemory(redis)
	rateLimiter := infra.NewRateLimiter(redis, 5, 2)

	// handler
	mux := http.NewServeMux()

	err = internal.NewApp(config, database, inMemory, rateLimiter, mux).
		AttachRoutes("/api/v1").
		Start()
	if err != nil {
		fmt.Printf("[ERROR] app start failed: %s", err)
	}
}
