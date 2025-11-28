package http

import (
	"github.com/gin-gonic/gin"

	"github.com/bielrodrigues/task-manager-pro-backend/internal/users"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")

	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// AUTH
	api.POST("/auth/register", users.RegisterHandler)
	api.POST("/auth/login", users.LoginHandler)
}
