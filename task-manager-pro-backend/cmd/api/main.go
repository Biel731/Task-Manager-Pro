package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/bielrodrigues/task-manager-pro-backend/internal/config"
	"github.com/bielrodrigues/task-manager-pro-backend/internal/database"
	"github.com/bielrodrigues/task-manager-pro-backend/internal/http"
	"github.com/bielrodrigues/task-manager-pro-backend/internal/users"
)

func main() {
	// Load env/config
	config.Load()

	// Connect to Postgres
	database.ConnectPostgres()

	// Connect to Redis
	database.ConnectRedis()

	// Migrate User
	users.Migrate()

	// Create Gin router
	r := gin.Default()

	// Register routes
	http.RegisterRoutes(r)

	log.Println("\r\n🚀 Server running at http://localhost:8080")
	r.Run(":8080")
}
