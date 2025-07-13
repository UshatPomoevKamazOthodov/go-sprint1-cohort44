package storage

import (
	"math/rand"
	"sync"
)

type URLPair struct {
	URL        string
	URLReduced string
}

var GlobalStorage *MemoryStorage

type MemoryStorage struct {
	mu       *sync.Mutex
	nextID   int
	urls     map[string]URLPair
	reversed map[string]string
}

func InitGlobalStorage() {
	GlobalStorage = &MemoryStorage{
		mu:       &sync.Mutex{},
		nextID:   1,
		urls:     make(map[string]URLPair),
		reversed: make(map[string]string),
	}
}

func (s *MemoryStorage) Save(original string) URLPair {
	s.mu.Lock()
	defer s.mu.Unlock()

	short := randomString(10)
	s.nextID++

	pair := URLPair{
		URL:        original,
		URLReduced: short,
	}

	s.urls[original] = pair
	s.reversed[short] = original

	return pair
}

func (s *MemoryStorage) GetByOriginal(original string) (URLPair, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	pair, found := s.urls[original]
	return pair, found
}

func (s *MemoryStorage) GetUrl(reduced string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	original, found := s.reversed[reduced]
	return original, found
}

func randomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
