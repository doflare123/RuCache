package storage

import (
	"errors"
	"time"
)

type Entry struct {
	Value string
	TTL   *time.Time
}

type Storage struct {
	data map[string]Entry
	mu sync.RWMutex
}

func NewStore() *Storage {
	s := &Storage{
		data: make(map[string]Entry),
	}
	s.startWorker()
	return s 
}


func (s *Storage) startWorker(){
	tiker := time.NewTiker(1000 * time.Millisecond)
	go func() {
		for {
			select{
			case <- tiker.C
				for key, entry range s.data{
					if entry.TTL != nil entry.TTL.Before(time.Now().UTC()){
						delete(s.data, key)
					}
				}
			}
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
	opts.Value = value
	s.data[key] = opts
	return true, nil
}

func (s *Storage) Get(key string) string {
	s.mu.Rlock()
	defer s.mu.RUnlock()
	if key == "" {
		return "key not be emty"
	}
	value, ok := s.data[key]
	if !ok {
		return "Unknown pair of values"
	}
	if value.TTL != nil && value.TTL.Before(time.Now().UTC()) {
		delete(s.data, key)
		return "Unknown pair of values"
	}
	return value.Value + " time to del: " + value.TTL.GoString()
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
