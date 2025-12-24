package http

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/bielrodrigues/task-manager-pro-backend/internal/ai"
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

	// ===== AI =====
	aiProvider := ai.NewOpenAIProvider(
		os.Getenv("OPENAI_API_KEY"),
		os.Getenv("OPENAI_MODEL"),
	)
	aiService := ai.NewService(aiProvider)
	aiHandler := ai.NewHandler(aiService)

	aiGroup := protected.Group("/ai")
	aiGroup.POST("/suggest-title", aiHandler.SuggestTitles)
	aiGroup.POST("/improve-description", aiHandler.ImproveDescription)
	aiGroup.POST("/generate-subtasks", aiHandler.GenerateSubtasks)
	// Resultado: POST /api/ai/suggest-title

	// ===== TASKS =====
	tasksGroup := protected.Group("/tasks")

	// LIST (sem cache, com filtros)
	tasksGroup.GET("", tasks.ListarTaskHandler)

	// SEARCH (com cache Redis) -> /api/tasks/search
	tasksGroup.GET("/search", tasks.SearchTasksHandler)

	// HISTORY -> /api/tasks/search/history  ✅ (corrigido)
	tasksGroup.GET("/search/history", tasks.GetSearchHistoryHandler)

	// GET por ID -> /api/tasks/:id
	tasksGroup.GET("/:id", tasks.GetTaskHandler)

	// CREATE -> /api/tasks
	tasksGroup.POST("", tasks.CreateTaskHandler)

	// UPDATE -> /api/tasks/:id
	tasksGroup.PUT("/:id", tasks.UpdateTaskHandler)

	// DELETE -> /api/tasks/:id
	tasksGroup.DELETE("/:id", tasks.DeleteTaskHandler)
}
