package storage

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) (*LocalStorage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, err
	}
	return &LocalStorage{basePath: basePath}, nil
}

func (s *LocalStorage) Save(path string, data []byte) error {
	fullPath := filepath.Join(s.basePath, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(fullPath, data, 0644)
}

func (s *LocalStorage) Read(path string) ([]byte, error) {
	fullPath := filepath.Join(s.basePath, path)
	return ioutil.ReadFile(fullPath)
}