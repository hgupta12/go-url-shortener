package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
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

type PostgresStorage struct{
	db *sql.DB
}

func NewPostgresStorage(host string, port int, user string, password string, dbname string) *PostgresStorage {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname))
	if err != nil {
		panic(err)
	}
	
	if err := db.Ping(); err != nil {
		panic(err)
	}

	query := `
		CREATE TABLE IF NOT EXISTS urls (
			hash VARCHAR(15) PRIMARY KEY,
			url TEXT NOT NULL
			);`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}

	return &PostgresStorage{
		db: db,
	}
}

func (s *PostgresStorage) Save(url string, hash string) error {
	query := `INSERT INTO urls (hash, url) VALUES ($1, $2)`
	_, err := s.db.Exec(query, hash, url)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStorage) Load(hash string) (string, error) {
	var url string
	query := `SELECT url FROM urls WHERE hash=$1`
	err := s.db.QueryRow(query, hash).Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}


type PostgresAndRedisStorage struct {
	postgres *PostgresStorage
	redis *RedisStorage
}

func NewPostgresAndRedisStorage(postgres *PostgresStorage, redis *RedisStorage) *PostgresAndRedisStorage {
	return &PostgresAndRedisStorage{
		postgres: postgres,
		redis: redis,
	}
}

func (s *PostgresAndRedisStorage) Save(url string, hash string) error {
	err := s.postgres.Save(url, hash)
	if err != nil {
		return err
	}
	_ = s.redis.Save(url, hash)
	
	return nil
}

func (s *PostgresAndRedisStorage) Load(hash string) (string, error) {
	url, err := s.redis.Load(hash)
	if err != nil {
		url, err = s.postgres.Load(hash)
		if err != nil {
			return "", err
		}
	}
	return url, nil
}
