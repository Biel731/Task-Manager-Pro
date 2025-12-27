package database

import (
	"log"

	"github.com/bielrodrigues/task-manager-pro-backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPostgres() {
	db, err := gorm.Open(postgres.Open(config.DatabaseUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("\r\nPostgreSQL connection failed:\r\n", err)
	}

	DB = db

	log.Println("Connected to PostgreSQL")
}
