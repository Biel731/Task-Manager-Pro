package database

import (
	"context"
	"log"

	"github.com/bielrodrigues/task-manager-pro-backend/internal/config"
	"github.com/redis/go-redis/v9"
)

var (
	Redis *redis.Client
	ctx   = context.Background()
)

func ConnectRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr: config.RedisURL,
	})

	_, err := Redis.Ping(ctx).Result()
	if err != nil {
		log.Fatal("\r\nRedis connection failed:", err)
	}

	log.Println("\r\nConnected to Redis")

}
