package redis

import (
	"github.com/idoyudha/eshop-product/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedis(cfg *config.Redis) *RedisClient {
	return &RedisClient{
		Client: redis.NewClient(RedisOptions(cfg)),
	}
}
