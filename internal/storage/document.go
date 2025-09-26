package storage

import (
	"os"
	"path/filepath"
)

type LocalFileStorage struct {
	BasePath string
}

func NewLocalFileStorage(basePath string) *LocalFileStorage {
	return &LocalFileStorage{BasePath: basePath}
}

func (s *LocalFileStorage) Save(fileName string, data []byte) (string, error) {
	path := filepath.Join(s.BasePath, fileName)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", err
	}
	return path, nil
}

func (s *LocalFileStorage) Delete(filePath string) error {
	if filePath == "" {
		return nil
	}
	return os.Remove(filePath)
}
