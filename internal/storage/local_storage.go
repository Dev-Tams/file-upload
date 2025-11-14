package storage

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/dev-tams/file-upload/internal/config"
)

type LocalStorage struct {
	BasePath string
}

func NewLocalStorage(cfg *config.AppConfig) (StorageProvider, error) {
	return &LocalStorage{
		BasePath: "uploads",
	}, nil
}

func (l *LocalStorage) Save(file *multipart.FileHeader, storedName string) (string, error) {
	dst := filepath.Join(l.BasePath, storedName)

	if err := os.MkdirAll(l.BasePath, 0755); err != nil {
		return "", err
	}

	if err := os.WriteFile(dst, nil, 0644); err != nil {
		return "", err
	}

	return dst, nil
}

func (l *LocalStorage) Delete(path string) error {
	return os.Remove(path)
}

func (l *LocalStorage) GetURL(path string) (string, error) {
	return fmt.Sprintf("/%s", path), nil
}
