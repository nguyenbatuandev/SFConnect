package service

import (
	"Product_Service/internal/cache"
	"time"
)

type RedisService struct {
	redisClient *cache.RedisCache
}

func NewRedisService(addr string, password string, db int) *RedisService {
	client := cache.NewRedisCache(addr, password, db)
	return &RedisService{
		redisClient: client,
	}
}

func (s *RedisService) Set(key string, dest interface{}, expiration time.Duration) error {
	return s.redisClient.Set(key, dest, expiration)
}

func (s *RedisService) Get(key string, dest interface{}) error {
	return s.redisClient.Get(key, dest)
}

func (s *RedisService) Delete(key string) error {
	return s.redisClient.Delete(key)
}

func (s *RedisService) Exists(key string) bool {
	exists := s.redisClient.Exists(key)
	return exists
}

func (s *RedisService) DeletePattern(pattern string) error {
	return s.redisClient.DeletePattern(pattern)
}
