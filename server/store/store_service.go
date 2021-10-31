package store

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type StorageService struct {
	redis *redis.Client
}

var (
	storeService = &StorageService{}
	ctx          = context.Background()
)

const CacheDuration = 6 * time.Hour

func InitializeStore(store *redis.Client) *StorageService {
	if store != nil {
		storeService.redis = store
		return storeService
	}

	addr, isAddrPresent := os.LookupEnv("REDIS_HOST")
	if !isAddrPresent {
		addr = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic("autschie")
	}

	storeService.redis = client
	return storeService
}

func SaveUrlMapping(shortUrl string, originalUrl string, userId string) {
	err := storeService.redis.Set(ctx, shortUrl, originalUrl, CacheDuration).Err()
	if err != nil {
		panic(err)
	}
}

func RetrieveOriginalUrl(shortUrl string) string {
	result, err := storeService.redis.Get(ctx, shortUrl).Result()
	if err != nil {
		panic(err)
	}

	return result
}
