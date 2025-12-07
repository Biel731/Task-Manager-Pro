package main

import (
	"github.com/gin-gonic/gin"

	"github.com/bielrodrigues/task-manager-pro-backend/internal/cache"
	"github.com/bielrodrigues/task-manager-pro-backend/internal/config"
	"github.com/bielrodrigues/task-manager-pro-backend/internal/database"
	"github.com/bielrodrigues/task-manager-pro-backend/internal/http"
	"github.com/bielrodrigues/task-manager-pro-backend/internal/tasks"
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
	tasks.Migrate()

	// Create Gin router
	r := gin.Default()

	// Register routes
	http.RegisterRoutes(r)

	// Redis
	redisClient := cache.NewClientRedis(config.RedisURL)
	tasksService := tasks.NewService(database.DB, redisClient)
}
