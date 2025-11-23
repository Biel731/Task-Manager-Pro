package users

import (
	"log"

	"github.com/bielrodrigues/task-manager-pro-backend/internal/database"
)

func Migrate() {
	err := database.DB.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("\r\n❌ Failed to migrate users table:", err)
	}

	log.Println("\r\n✅ Users table migrated")
}
