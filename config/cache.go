package config

import (
	"fmt"
	"github.com/caarlos0/env/v8"
)

type RedisConfig struct {
	Host     string `env:"REDIS_HOST"`
	Port     int    `env:"REDIS_PORT"`
	Address  string
	Password string `env:"REDIS_PASSWORD"`
}

type CacheConfig struct {
	Type  string `env:"CACHE_TYPE"`
	Redis *RedisConfig
}

var redisConfig *RedisConfig
var cacheConfig *CacheConfig

func GetRedis() (*RedisConfig, error) {
	if redisConfig != nil {
		return redisConfig, nil
	}
	redisConfig = &RedisConfig{}
	if err := env.Parse(redisConfig); err != nil {
		return nil, err
	}
	if redisConfig.Host == "" {
		redisConfig.Host = "localhost"
	}
	if redisConfig.Port == 0 {
		redisConfig.Port = 6379
	}

	redisConfig.Address = fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)
	return redisConfig, nil
}

func GetCache() (*CacheConfig, error) {
	if cacheConfig != nil {
		return cacheConfig, nil
	}
	redisConfig, err := GetRedis()
	if err != nil {
		return nil, err
	}
	cacheConfig = &CacheConfig{
		Redis: redisConfig,
	}
	if err := env.Parse(cacheConfig); err != nil {
		return nil, err
	}

	return cacheConfig, nil
}
