package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisClient struct {
	Client *redis.Client
}

func NewClientRedis(addr string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisClient{Client: rdb}
}

func (r *RedisClient) Set(key string, value string, ttl time.Duration) error {
	return r.Client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisClient) LPush(key string, value string) error {
	return r.Client.LPush(ctx, key, value).Err()
}

func (r *RedisClient) LTrim(key string, start, stop int64) error {
	return r.Client.LPush(ctx, key, start, stop).Err()
}
