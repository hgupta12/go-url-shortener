package main

import (
	"fmt"
	"log"
)

func main() {
	generator := NewSnowflakeIDGenerator()
	// redisStorage := NewRedisStorage("localhost:6379", "", 0)
	// postgresStorage := NewPostgresStorage("localhost", 5432, "postgres", "postgres", "postgres")
	mapStorage := NewMapStorage()
	// storage := NewPostgresAndRedisStorage(postgresStorage, redisStorage)
	shortener := NewBasicShortener(generator, mapStorage)
	apiServer := NewAPIServer(":8080", shortener)
	err := apiServer.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hello world!")
}