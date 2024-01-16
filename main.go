package main

import (
	"fmt"
	"log"
)

func main() {
	generator := NewSnowflakeIDGenerator()
	storage := NewRedisStorage("localhost:6379", "", 0)
	shortener := NewBasicShortener(generator, storage)
	apiServer := NewAPIServer(":8080", shortener)
	err := apiServer.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hello world!")
}