package storage

import (
	"encoding/gob"
	"os"
	"path/filepath"
)

func (s *Storage) SaveDataToFile() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	path, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(path, dirName), os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(path, dirName, snapshotFileName))
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(s.data); err != nil {
		return err
	}
	return nil
}

func (s *Storage) LoadDataFromFile() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	path, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(path, dirName), os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.Open(filepath.Join(path, dirName, snapshotFileName))
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&s.data); err != nil {
		return err
	}
	return nil
}
