package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)

type Redis[T any] struct {
	Interface[T]
	redis             *redis.Client
	context           context.Context
	marshalCallback   func(any) ([]byte, error)
	unmarshalCallback func([]byte, any) error
}

func NewRedis[T any](options *redis.Options) Redis[T] {
	return Redis[T]{
		redis:             redis.NewClient(options),
		context:           context.Background(),
		marshalCallback:   json.Marshal,
		unmarshalCallback: json.Unmarshal,
	}
}

func (r Redis[T]) Set(key string, value T, time time.Duration) {
	r.redis.Set(r.context, key, value, time)
}

func (r Redis[T]) Remember(key string, time time.Duration, callback func() T) T {
	if v := r.Get(key); v != nil {
		return v
	}
	value := callback()
	r.Set(key, value, time)
	return value
}

func (r Redis[T]) Get(key string) T {
	return r.redis.Get(r.context, key)
}
