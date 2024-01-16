package main

type Request struct {
	Url string `json:"url"`
}

type Response struct {
	Url string `json:"url"`
}

type HashGenerator interface {
	Generate() string
}

type Shortener interface {
	Shorten(url string) string
	Resolve(url string) (string, error)
}

type Storage interface {
	Save(url string, hash string) error
	Load(hash string) (string, error)
}