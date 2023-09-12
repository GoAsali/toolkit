package services

import (
	"context"
	"encoding/json"
	"errors"
	gocache "github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	"github.com/goasali/toolkit/cache"
	"github.com/goasali/toolkit/global"
	"strconv"
	"time"
)

type Temp struct {
	ctx   context.Context
	cache *gocache.Cache[string]
}

type TemporaryLinksOption struct {
	userId     uint
	expiration time.Duration
}

type TemporaryLinksOptionFunc func(*TemporaryLinksOption)

func WithUser(user uint) TemporaryLinksOptionFunc {
	return func(option *TemporaryLinksOption) {
		option.userId = user
	}
}

func WithExpiration(expiration time.Duration) TemporaryLinksOptionFunc {
	return func(option *TemporaryLinksOption) {
		option.expiration = expiration
	}
}

func NewTemporary() (*Temp, error) {
	cacheInstance, err := cache.New[string]()

	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	return &Temp{
		ctx,
		cacheInstance,
	}, nil
}

func (temp *Temp) GenerateFileLink(path string, functions ...TemporaryLinksOptionFunc) error {
	option := getOption(functions)

	value := make(map[string]string)
	value["path"] = path

	key := path

	if option.userId != 0 {
		idStr := strconv.Itoa(int(option.userId))
		value["user"] = idStr
		key += idStr
	}

	key = global.GetMD5Hash(key)
	stringify, err := json.Marshal(value)
	if err != nil {
		return err
	}
	du, _ := time.ParseDuration("1m")
	return temp.cache.Set(temp.ctx, key, string(stringify), store.WithExpiration(du))
}

func (temp *Temp) GetFileLink(key string, functions ...TemporaryLinksOptionFunc) (string, error) {
	option := getOption(functions)
	valueString, err := temp.cache.Get(temp.ctx, key)
	if err != nil {
		if errors.Is(err, &store.NotFound{}) {
			return "", nil
		}
		return "", err
	}

	var value map[string]string
	if err := json.Unmarshal([]byte(valueString), &value); err != nil {
		return "", err
	}

	if _, ok := value["user"]; option.userId != 0 && ok {
		userId, err := strconv.Atoi(value["user"])
		if err != nil {
			return "", err
		}
		if userId != int(option.userId) {
			return "", nil
		}
	}
	path, ok := value["path"]
	if !ok {
		return "", nil
	}
	return path, nil
}

func getOption(functions []TemporaryLinksOptionFunc) TemporaryLinksOption {
	option := TemporaryLinksOption{}
	for _, optionFunc := range functions {
		optionFunc(&option)
	}
	return option
}
