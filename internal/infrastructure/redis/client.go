package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/vlladoff/micro-learn/internal/config"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	addr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
	
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return rdb
}
