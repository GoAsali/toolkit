package cache

import (
	redisstore "github.com/eko/gocache/store/redis/v4"
	"github.com/goasali/toolkit/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisStore() (*redisstore.RedisStore, error) {
	redisConfig, err := config.GetRedis()
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Address,
		Password: redisConfig.Password,
	})

	return redisstore.NewRedis(client), nil
}
