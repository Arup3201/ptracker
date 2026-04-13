package main

import (
	"fmt"
	"log"

	"github.com/ptracker/cmd/app"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

	err = app.NewApp("/api/v1", config, db).Start()
	if err != nil {
		fmt.Printf("[ERROR] app start: %s", err)
	}
}
