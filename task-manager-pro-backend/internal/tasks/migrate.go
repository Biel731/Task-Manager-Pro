package tasks

import (
	"log"

	"github.com/bielrodrigues/task-manager-pro-backend/internal/database"
)

func Migrate() {
	err := database.DB.AutoMigrate(&Tag{}, &Task{})
	if err != nil {
		log.Fatal("❌ Failed to migrate tasks/tags tables:", err)
	}

	log.Println("✅ Tasks & Tags tables migrated")

}
