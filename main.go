package main

import (
	"fmt"
	"log"
)

func main() {
	numberGenerator := NewNumberGenerator()
	redisStorage := NewRedisStorage("localhost:6379", "", 0)
	shortener := NewBasicShortener(numberGenerator, redisStorage)
	apiServer := NewAPIServer(":8080", shortener)
	err := apiServer.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hello world!")
}