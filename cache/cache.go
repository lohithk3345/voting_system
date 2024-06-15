package cache

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/lohithk3345/voting_system/config"
	"github.com/lohithk3345/voting_system/types"
	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	client *redis.Client
	mu     sync.Mutex
}

func NewCacheService() *CacheService {
	return &CacheService{
		client: redis.NewClient(&redis.Options{
			Addr:     config.EnvMap["REDIS_URL"],
			Password: "",
			DB:       0,
		}),
		mu: sync.Mutex{},
	}
}

func (cache *CacheService) Set(key string, value interface{}) error {
	ctx := context.Background()

	cache.mu.Lock()

	err := cache.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		log.Println("Error Setting Value", err)
		return CacheError{Code: SET_ERROR, Message: err.Error()}
	}

	cache.mu.Unlock()

	return nil
}

func (cache *CacheService) SetTokens(key string, value types.Tokens) error {
	ctx := context.Background()

	input_key := "user:" + key + "tokens"

	marshaledToken, errMarshal := json.Marshal(value)
	if errMarshal != nil {
		log.Println("Error Setting Tokens", errMarshal)
		return CacheError{Code: SET_MARSHAL_ERROR, Message: errMarshal.Error()}
	}

	cache.mu.Lock()

	err := cache.client.Set(ctx, input_key, marshaledToken, 0).Err()
	if err != nil {
		log.Println("Error Setting Tokens", err)
		return CacheError{Code: SET_ERROR, Message: err.Error()}
	}

	cache.mu.Unlock()

	return nil
}

func (cache *CacheService) GetTokens(key string) (*types.Tokens, error) {
	ctx := context.Background()

	cache.mu.Lock()

	input_key := "user:" + key + "tokens"

	value, err := cache.client.Get(ctx, input_key).Result()
	if err != nil {
		log.Println("Error Getting Value", err)
		return nil, CacheError{Code: GET_ERROR, Message: err.Error()}
	}

	cache.mu.Unlock()

	var tokens types.Tokens

	json.Unmarshal([]byte(value), &tokens)

	return &tokens, nil
}

func (cache *CacheService) Get(key string) (interface{}, error) {
	ctx := context.Background()

	cache.mu.Lock()

	value, err := cache.client.Get(ctx, key).Result()
	if err != nil {
		log.Println("Error Getting Value", err)
		return nil, CacheError{Code: GET_ERROR, Message: err.Error()}
	}

	cache.mu.Unlock()

	return value, nil
}

func (cache *CacheService) DeleteTokenKey(key string) error {
	ctx := context.Background()

	input_key := "user:" + key + "tokens"

	cache.mu.Lock()

	err := cache.client.Del(ctx, input_key).Err()
	if err != nil {
		log.Println("Error Deleting Key", err)
		return CacheError{Code: SET_ERROR, Message: err.Error()}
	}

	cache.mu.Unlock()

	return nil
}
