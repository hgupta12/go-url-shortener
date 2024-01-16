package main

import (
	"fmt"
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