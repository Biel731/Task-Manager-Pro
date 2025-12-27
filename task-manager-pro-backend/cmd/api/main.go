package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/bielrodrigues/task-manager-pro-backend/internal/cache"
	"github.com/bielrodrigues/task-manager-pro-backend/internal/config"
	"github.com/bielrodrigues/task-manager-pro-backend/internal/database"
	internalhttp "github.com/bielrodrigues/task-manager-pro-backend/internal/http"
	"github.com/bielrodrigues/task-manager-pro-backend/internal/tasks"
	"github.com/bielrodrigues/task-manager-pro-backend/internal/users"
)

func main() {
	// Carrega variáveis de ambiente / config
	config.Load()

	// Conecta no Postgres
	database.ConnectPostgres()

	redisClient := cache.NewClientRedis(config.RedisURL)
	tasks.SetRedisClient(redisClient)

	// Migrations
	users.Migrate()
	tasks.Migrate()

	// Cria router Gin
	r := gin.Default()

	// Registra as rotas (auth, users, tasks, etc.)
	internalhttp.RegisterRoutes(r)

	// Sobe o servidor
	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}
