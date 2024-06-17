package cache

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"

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

func (cache *CacheService) InitWSManager() ([]*types.InitVoteData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	roomKeys, err := cache.client.Keys(ctx, "voting:*").Result()
	if err != nil {
		return nil, err
	}

	// var voteData []*types.VoteData
	// var keys []string

	var initData []*types.InitVoteData

	for _, roomKey := range roomKeys {
		roomID := strings.TrimPrefix(roomKey, "voting:")
		data, err := cache.GetVoteDataByRoomId(roomID)
		if err != nil {
			return nil, err
		}

		initData = append(initData, &types.InitVoteData{
			Data: data,
			Keys: roomID,
		})
	}

	return initData, nil
}

func (cache *CacheService) Set(key string, value interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

func (cache *CacheService) CreateRoom(roomId string, data *types.VoteData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cache.mu.Lock()
	input_key := "voting:" + roomId

	marshalData, errMarshal := json.Marshal(data)
	if errMarshal != nil {
		log.Println("Error Setting Tokens", errMarshal)
		return CacheError{Code: SET_MARSHAL_ERROR, Message: errMarshal.Error()}
	}

	err := cache.client.Set(ctx, input_key, marshalData, 0).Err()

	if err != nil {
		log.Println("Error Setting Tokens", err)
		return CacheError{Code: SET_ERROR, Message: err.Error()}
	}

	cache.mu.Unlock()

	return nil
}

func (cache *CacheService) SetVoteByRoomId(roomId string, option string) (*types.VoteData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cache.mu.Lock()

	input_key := "voting:" + roomId

	value, err := cache.client.Get(ctx, input_key).Result()
	if err != nil {
		log.Println("Error Getting Value", err)
		return nil, CacheError{Code: GET_ERROR, Message: err.Error()}
	}

	// cache.mu.Unlock()

	var voteData types.VoteData

	json.Unmarshal([]byte(value), &voteData)

	voteData.Options[option]++

	marshalData, errMarshal := json.Marshal(voteData)
	if errMarshal != nil {
		log.Println("Error Setting Tokens", errMarshal)
		return nil, CacheError{Code: SET_MARSHAL_ERROR, Message: errMarshal.Error()}
	}

	// cache.mu.Lock()

	errSet := cache.client.Set(ctx, input_key, marshalData, 0).Err()
	if errSet != nil {
		log.Println("Error Setting Tokens", err)
		return nil, CacheError{Code: SET_ERROR, Message: errSet.Error()}
	}

	cache.mu.Unlock()

	return &voteData, nil
}

func (cache *CacheService) GetVoteDataByRoomId(roomId string) (*types.VoteData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cache.mu.Lock()

	input_key := "voting:" + roomId

	value, err := cache.client.Get(ctx, input_key).Result()
	if err != nil {
		log.Println("Error Getting Value", err)
		return nil, CacheError{Code: GET_ERROR, Message: err.Error()}
	}

	cache.mu.Unlock()

	var voteData types.VoteData

	json.Unmarshal([]byte(value), &voteData)

	return &voteData, nil
}
