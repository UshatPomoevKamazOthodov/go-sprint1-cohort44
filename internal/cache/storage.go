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
	sync.RWMutex
	urls     map[string]URLPair // ключ = оригинальный URL
	reversed map[string]string  // ключ = сокращённый URL
	nextID   int                // для генерации ID
}

func InitGlobalStorage() {
	GlobalStorage = &MemoryStorage{
		urls:     make(map[string]URLPair),
		reversed: make(map[string]string),
		nextID:   1,
	}
}

func (s *MemoryStorage) Save(original string) URLPair {
	s.Lock()
	defer s.Unlock()

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
	s.RLock()
	defer s.RUnlock()

	pair, found := s.urls[original]
	return pair, found
}

func (s *MemoryStorage) GetUrl(reduced string) (string, bool) {
	s.RLock()
	defer s.RUnlock()

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
