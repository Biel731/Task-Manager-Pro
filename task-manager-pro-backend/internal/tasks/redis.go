package tasks

import "github.com/bielrodrigues/task-manager-pro-backend/internal/cache"

var redisClient *cache.RedisClient

func SetRedisClient(client *cache.RedisClient) {
	redisClient = client
}
