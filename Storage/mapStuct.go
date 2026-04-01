package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	snapshotFileName = "snapshot.gob"
	dirName          = "RuCache"
)

type EntryType string

const (
	TypeString EntryType = "string"
	TypeHash   EntryType = "hash"
)

type Entry struct {
	Type     EntryType
	StrValue string
	Value    map[string]string
	TTL      *time.Time
}

type Storage struct {
	data map[string]Entry
	mu   sync.RWMutex
}

func NewStore() *Storage {
	path, _ := os.UserCacheDir()
	s := &Storage{
		data: make(map[string]Entry),
	}
	if _, err := os.Stat(filepath.Join(path, dirName, snapshotFileName)); err == nil {
		s.LoadDataFromFile()
	}
	s.startWorker()
	return s
}

func (s *Storage) startWorker() {
	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			s.mu.Lock()
			for key, entry := range s.data {
				if entry.TTL != nil && entry.TTL.Before(time.Now().UTC()) {
					fmt.Printf("Key %s expired and removed\n", key)
					delete(s.data, key)
				}
			}
			s.mu.Unlock()
		}
	}()
}

func (s *Storage) Set(key string, value string, ttl *time.Duration) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if key == "" || value == "" {
		return false, errors.New("key or value not be emty")
	}
	var opts Entry
	if ttl != nil {
		expAt := time.Now().UTC().Add(*ttl)
		opts.TTL = &expAt
	}
	if existing, ok := s.data[key]; ok {
		if existing.Type != TypeString {
			return false, errors.New("key exists with different type")
		}
	}
	opts.Type = TypeString
	opts.StrValue = value
	s.data[key] = opts
	return true, nil
}

func (s *Storage) HSet(key string, field [][]string, ttl *time.Duration) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if key == "" || len(field) == 0 {
		return false, errors.New("key or fields not be emty")
	}
	var opts Entry
	if ttl != nil {
		expAt := time.Now().UTC().Add(*ttl)
		opts.TTL = &expAt
	}
	opts.Type = TypeHash
	if existing, ok := s.data[key]; ok {
		if existing.Type != TypeHash {
			return false, errors.New("key exists with different type")
		}
		opts.Value = existing.Value
	} else {
		opts.Value = make(map[string]string)
	}
	for _, pair := range field {
		if len(pair) >= 2 {
			opts.Value[pair[0]] = pair[1]
		}
	}
	s.data[key] = opts
	return true, nil
}

func (s *Storage) HGet(key string, field string) (*string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if key == "" || field == "" {
		return nil, errors.New("key or field not be emty")
	}
	value, err := s.Cheсker(key, field)
	if err != nil {
		return nil, err
	}
	fieldValue, ok := value.Value[field]
	if !ok {
		return nil, errors.New("Unknown pair of values")
	}
	return &fieldValue, nil
}

func (s *Storage) HGetAll(key string) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, err := s.CheсkerKey(key)
	if err != nil {
		return nil, err
	}
	return value.Value, nil
}

func (s *Storage) Get(key string) (*string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if key == "" {
		return nil, errors.New("key not be emty")
	}
	value, ok := s.data[key]
	if !ok {
		return nil, errors.New("Unknown pair of values")
	}
	if value.TTL != nil && value.TTL.Before(time.Now().UTC()) {
		delete(s.data, key)
		return nil, errors.New("Unknown pair of values")
	}
	return &value.StrValue, nil
}

func (s *Storage) Del(key string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if key == "" {
		return false, errors.New("key not be emty")
	}
	delete(s.data, key)
	return true, nil
}

func (s *Storage) HDel(key string, field string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.Cheсker(key, field)
	if err != nil {
		return false, err
	}
	delete(s.data[key].Value, field)
	return true, nil
}

func (s *Storage) Cheсker(key string, field string) (*Entry, error) {
	if key == "" || field == "" {
		return nil, errors.New("key or field not be emty")
	}
	value, ok := s.data[key]
	if !ok {
		return nil, errors.New("Unknown pair of values")
	}
	if value.Type != TypeHash {
		return nil, errors.New("key exists with different type")
	}
	if value.TTL != nil && value.TTL.Before(time.Now().UTC()) {
		delete(s.data, key)
		return nil, errors.New("Unknown pair of values")
	}
	return &value, nil
}

func (s *Storage) CheсkerKey(key string) (*Entry, error) {
	if key == "" {
		return nil, errors.New("key or field not be emty")
	}
	value, ok := s.data[key]
	if !ok {
		return nil, errors.New("Unknown pair of values")
	}
	if value.Type != TypeHash {
		return nil, errors.New("key exists with different type")
	}
	if value.TTL != nil && value.TTL.Before(time.Now().UTC()) {
		delete(s.data, key)
		return nil, errors.New("Unknown pair of values")
	}
	return &value, nil
}
