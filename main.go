package main

import (
	"fmt"
	"log"
)

func main() {
	shortener := NewBasicShortener()
	apiServer := NewAPIServer(":8080", shortener)
	err := apiServer.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hello world!")
}