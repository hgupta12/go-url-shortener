package main

type BasicShortener struct {
	hashGenerator HashGenerator
	storage       Storage
}

func NewBasicShortener(hashGenerator HashGenerator, storage Storage) *BasicShortener {
	return &BasicShortener{
		hashGenerator: hashGenerator,
		storage: storage,
	}
}

func (s *BasicShortener) Shorten(url string) string {
	hash := s.hashGenerator.Generate()
	s.storage.Save(url, hash)
	return hash
}

func (s *BasicShortener) Resolve(hash string) (string, error) {
	url, err := s.storage.Load(hash)
	if err != nil {
		return "", err
	}

	return url, nil
}