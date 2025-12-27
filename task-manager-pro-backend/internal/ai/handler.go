package ai

import (
	"log"
	"net/http"
	"time"

	"github.com/bielrodrigues/task-manager-pro-backend/internal/auth"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SuggestTitles(c *gin.Context) {
	start := time.Now()

	userID, _ := auth.GetUserID(c)

	var req SuggestTitleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[AI] user=%d endpoint=/ai/suggest-title status=400 ms=%d err=bind_json",
			userID, time.Since(start).Milliseconds(),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	if msg := req.Validate(); msg != "" {
		log.Printf("[AI] user=%d endpoint=/ai/suggest-title status=400 ms=%d err=%s",
			userID, time.Since(start).Milliseconds(), msg,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	titles, err := h.service.SuggestTitles(c.Request.Context(), req.Description)
	if err != nil {
		log.Printf("[AI] user=%d endpoint=/ai/suggest-title status=500 ms=%d err=%v",
			userID, time.Since(start).Milliseconds(), err,
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate titles"})
		return
	}

	log.Printf("[AI] user=%d endpoint=/ai/suggest-title status=200 ms=%d",
		userID, time.Since(start).Milliseconds(),
	)

	c.JSON(http.StatusOK, SuggestTitleResponse{Titles: titles})
}

func (h *Handler) ImproveDescription(c *gin.Context) {
	start := time.Now()
	userID, _ := auth.GetUserID(c)

	var req ImproveDescriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[AI] user=%d endpoint=/ai/improve-description status=400 ms=%d err=bind_json",
			userID, time.Since(start).Milliseconds(),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	if msg := req.Validate(); msg != "" {
		log.Printf("[AI] user=%d endpoint=/ai/improve-description status=400 ms=%d err=%s",
			userID, time.Since(start).Milliseconds(), msg,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	improved, bullets, err := h.service.ImproveDescription(c.Request.Context(), req.Title, req.Description)
	if err != nil {
		log.Printf("[AI] user=%d endpoint=/ai/improve-description status=500 ms=%d err=%v",
			userID, time.Since(start).Milliseconds(), err,
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to improve description"})
		return
	}

	log.Printf("[AI] user=%d endpoint=/ai/improve-description status=200 ms=%d",
		userID, time.Since(start).Milliseconds(),
	)

	c.JSON(http.StatusOK, ImproveDescriptionResponse{
		ImprovedDescription: improved,
		Bullets:             bullets,
	})
}
