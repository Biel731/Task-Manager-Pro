package http

import (
	"github.com/gin-gonic/gin"

	"github.com/bielrodrigues/task-manager-pro-backend/internal/auth"
	"github.com/bielrodrigues/task-manager-pro-backend/internal/tasks"
	"github.com/bielrodrigues/task-manager-pro-backend/internal/users"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")

	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// AUTH (público)
	api.POST("/auth/register", users.RegisterHandler)
	api.POST("/auth/login", users.LoginHandler)

	// Rotas protegidas
	protected := api.Group("/")
	protected.Use(auth.AuthMiddleware())

	// TASKS
	protected.GET("/tasks", tasks.ListarTaskHandler)
	protected.GET("/tasks/:id", tasks.GetTaskHandler)
	protected.POST("/tasks", tasks.CreateTaskHandler)
	protected.PUT("/tasks/:id", tasks.UpdateTaskHandler)
	protected.DELETE("/tasks/:id", tasks.DeleteTaskHandler)
}
