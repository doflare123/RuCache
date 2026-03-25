package storage

import (
	"errors"
	"time"
)

type Entry struct {
	Value string
	TTL   time.Time
}

type Storage struct {
	data map[string]Entry
}

func NewStore() *Storage {
	return &Storage{
		data: make(map[string]Entry),
	}
}

func (s *Storage) Set(key string, value string) (bool, error) {
	if key == "" || value == "" {
		return false, errors.New("key or value not be emty")
	}
	s.data[key] = Entry{value, time.Time{}}
	return true, nil
}

func (s *Storage) Get(key string) (string, error) {
	if key == "" {
		return "", errors.New("key not be emty")
	}
	value, ok := s.data[key]
	if !ok {
		return "", errors.New("Unknown pair of values")
	}
	return value.Value, nil
}

func (s *Storage) Del(key string) (bool, error) {
	if key == "" {
		return false, errors.New("key not be emty")
	}
	delete(s.data, key)
	return true, nil
}
