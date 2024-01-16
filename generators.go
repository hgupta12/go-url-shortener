package main

import (
	"fmt"
	"sync"
)

type NumberGenerator struct{
	counter int
	mutex sync.Mutex
}

func NewNumberGenerator() *NumberGenerator {
	return &NumberGenerator{}
}

func (g *NumberGenerator) Generate(url string) string {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.counter++
	return fmt.Sprintf("%d", g.counter)
}