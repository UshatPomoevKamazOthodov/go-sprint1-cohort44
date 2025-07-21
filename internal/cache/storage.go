package storage

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type URLPair struct {
	UUID       string `json:"uuid"`
	URL        string `json:"url"`
	URLReduced string `json:"url_reduced"`
}

var GlobalStorage *MemoryStorage

type MemoryStorage struct {
	mu          *sync.Mutex
	storagePath string
	nextID      int
	urls        map[string]URLPair
	reversed    map[string]URLPair
}

func InitGlobalStorage(path string) *MemoryStorage {
	ms := &MemoryStorage{
		mu:          &sync.Mutex{},
		storagePath: path,
		urls:        make(map[string]URLPair),
		reversed:    make(map[string]URLPair),
		nextID:      1,
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := ms.createEmptyFile(); err != nil {
			log.Printf("Failed to create storage file: %v", err)
		}
	}

	if _, err := ms.LoadFromFile(); err != nil {
		log.Printf("Failed to load from file: %v", err)
	}

	return ms
}

func (s *MemoryStorage) LoadFromFile() ([]URLPair, error) {
	os.Open(s.storagePath)
	data, err := os.ReadFile(s.storagePath)
	if os.IsNotExist(err) {
		return nil, errors.Wrap(os.ErrNotExist, "Storage file not found")
	}
	if err != nil {
		return nil, errors.Wrap(err, "read file")
	}

	var records []URLPair
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}

	for _, record := range records {
		id, err := strconv.Atoi(record.UUID)
		if err != nil {
			log.Printf("Invalid UUID: %s", record.UUID)
			continue
		}

		s.urls[record.URL] = record
		s.reversed[record.URLReduced] = record

		if id >= s.nextID {
			s.nextID = id + 1
		}
	}

	return records, nil
}

func (s *MemoryStorage) SaveToFile() error {
	records := make([]URLPair, 0, len(s.urls))
	for _, record := range s.urls {
		records = append(records, record)
	}

	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return errors.Wrap(err, "marshal")
	}

	return os.WriteFile(s.storagePath, data, 0644)
}

func (s *MemoryStorage) Save(original string) string {
	uuid := strconv.Itoa(s.nextID)
	short := randomString(10)
	s.nextID++

	record := URLPair{
		UUID:       uuid,
		URLReduced: short,
		URL:        original,
	}
	s.urls[record.URL] = record
	s.reversed[record.URLReduced] = record

	if err := s.SaveToFile(); err != nil {
		log.Printf("Failed to save to file: %v", err)
	}

	return short
}

func (s *MemoryStorage) GetByOriginal(original string) (URLPair, bool) {
	if s.urls != nil && len(s.urls) > 0 {
		pair, found := s.urls[original]
		return pair, found
	}
	return URLPair{}, false
}

func (s *MemoryStorage) GetUrl(reduced string) (string, bool) {
	record, found := s.reversed[reduced]
	log.Println(record, found)
	if !found {
		return "", false
	}
	return record.URL, true
}

func (s *MemoryStorage) createEmptyFile() error {
	dir := filepath.Dir(s.storagePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.Wrap(err, "failed to create directory")
	}

	file, err := os.Create(s.storagePath)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	defer file.Close()

	_, err = file.WriteString("[]")
	return err
}

func randomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
