package tasks

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bielrodrigues/task-manager-pro-backend/internal/cache"
	"gorm.io/gorm"
)

type Service struct {
	db    *gorm.DB
	cache *cache.RedisClient
}

func NewService(db *gorm.DB, cache *cache.RedisClient) *Service {
	return &Service{db: db, cache: cache}
}

func (s *Service) SearchTasks(ctx context.Context, userID uint, query string) ([]Task, error) {

	// Montar as chaves
	queryHash := hashQuery(query)
	cacheKey := "tmpro:search:result:" + string(rune(userID)) + ":" + queryHash
	historyKey := "tmpro:search:history:" + string(rune(userID))

	// Pegar informações do cache
	if s.cache != nil {
		if cached, err := s.cache.Get(cacheKey); err == nil && cached != "" {
			var tasks []Task
			if err := json.Unmarshal([]byte(cached), &tasks); err == nil {
				return tasks, nil
			}
		}
	}

	// Se não tem cache, busca no banco
	var task []Task

	// busca por título ou descrição
	err := s.db.WithContext(ctx).
		Where("user_id = ? AND (title ILIKE ? OR description ILIKE ?)", userID, "%"+query+"%", "%"+query+"%").
		Find(&task).Error

	if err != nil {
		return nil, err
	}

	// Salva o resultado do cache
	if s.cache != nil {
		b, _ := json.Marshal(task)
		_ = s.cache.Set(cacheKey, string(b), 30*time.Second)
		_ = s.cache.LPush(historyKey, query)
		_ = s.cache.LTrim(historyKey, 0, 9)
	}

	return task, nil
}
