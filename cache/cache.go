package cache

import (
	"errors"
	gocache "github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	"github.com/goasali/toolkit/config"
	log "github.com/sirupsen/logrus"
)

var cacheStore store.StoreInterface

func NewStore() (store.StoreInterface, error) {
	if cacheStore != nil {
		return cacheStore, nil
	}
	cacheConfig, err := config.GetCache()
	if err != nil {
		log.Errorf("NewStore config error: %v", err)
		return nil, err
	}
	if cacheConfig.Type == "redis" {
		cacheStore, err = NewRedisStore()
		if err != nil {
			return nil, err
		}
		return cacheStore, nil
	}
	return nil, errors.New("could not create cache instance")
}

func New[T any]() (*gocache.Cache[T], error) {
	cacheStore, err := NewStore()
	if err != nil {
		return nil, err
	}
	cacheInstance := gocache.New[T](cacheStore)
	return cacheInstance, err
}
