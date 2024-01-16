package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type MapStorage struct {
	store map[string]string
}

func NewMapStorage() *MapStorage {
	return &MapStorage{
		store: make(map[string]string),
	}
}

func (s *MapStorage) Save(url string, hash string) error {
	s.store[hash] = url
	return nil
}

func (s *MapStorage) Load(hash string) (string, error) {
	url, ok := s.store[hash]
	if !ok {
		return "", fmt.Errorf("URL not found")		
	}
	return url, nil
}

type RedisStorage struct {
	client *redis.Client
	ctx context.Context
}

func NewRedisStorage(address string, password string, db int) *RedisStorage {
	client := redis.NewClient(&redis.Options{
		Addr: address,
		Password: password,
		DB: db,
	})
	ctx := context.Background()
	return &RedisStorage{
		client: client,
		ctx: ctx,
	}
}

func (s *RedisStorage) Save(url string, hash string) error {
	err := s.client.Set(s.ctx, "url:" + hash, url, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *RedisStorage) Load(hash string) (string, error) {
	url, err := s.client.Get(s.ctx, "url:" + hash).Result()
	if err != nil {
		return "", err
	}
	return url, nil
}