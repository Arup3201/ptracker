package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"github.com/ptracker/cmd/app"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	FRONTEND_VERIFY_URL = "http://localhost:5173/verify"
)

func readRSAPrivateKey(filename string) (*rsa.PrivateKey, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	block, _ := pem.Decode(bytes)
	parseResult, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("x509 parse pkce#1 private key: %w", err)
	}

	privateKey := parseResult.(*rsa.PrivateKey)
	return privateKey, nil
}

func main() {
	var err error

	config := &app.Config{}
	err = config.LoadFromEnv()
	if err != nil {
		log.Fatalf("[ERROR] config load from env: %s", err)
	}

	db, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", config.DBHost, config.DBPort,
		config.DBUser, config.DBPass, config.DBName)),
		&gorm.Config{})
	if err != nil {
		log.Fatalf("[ERROR] gorm open: %s", err)
	}

	app.Migrate(db)

	redis := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	privateKey, err := readRSAPrivateKey("private.key")
	if err != nil {
		log.Fatalf("rsa generate key: %s\n", err)
	}

	app := app.NewApp(
		"/api/v1",
		config,
		db,
		redis,
		privateKey,
		FRONTEND_VERIFY_URL,
	)
	app.AllowedCrossOrigins = []string{"http://localhost:5173"}
	err = app.Start()
	if err != nil {
		fmt.Printf("[ERROR] app start: %s", err)
	}
}
