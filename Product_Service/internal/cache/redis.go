package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(addr string, password string, db int) *RedisCache {
	log.Printf("Connecting to Redis at %s, DB %d", addr, db)

	// Nếu password rỗng, không gửi password
	var redisPassword string
	if password != "" {
		redisPassword = password
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisPassword,
		DB:       db,
	})

	ctx := context.Background()

	return &RedisCache{
		client: client,
		ctx:    ctx,
	}
}

func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = r.client.Set(r.ctx, key, jsonData, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisCache) Get(key string, dest interface{}) error {
	val, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return err
	}

	// Giải mã JSON lưu trong redis vào biến dest
	return json.Unmarshal([]byte(val), dest)
}
func (r *RedisCache) Delete(key string) error {
	_, err := r.client.Del(r.ctx, key).Result()
	return err
}

func (r *RedisCache) Exists(key string) bool {
	exists, err := r.client.Exists(r.ctx, key).Result()
	if err != nil {
		return false
	}
	return exists > 0
}

func (r *RedisCache) DeletePattern(pattern string) error {
	var cursor uint64
	for {
		keys, nextCursor, err := r.client.Scan(r.ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := r.client.Del(r.ctx, keys...).Err(); err != nil {
				return err
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}
