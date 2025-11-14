package storage

import (
    "mime/multipart"
    "github.com/dev-tams/file-upload/internal/config"
)

type StorageProvider interface {
    Save(file *multipart.FileHeader, storedName string) (string, error)
    Delete(path string) error
    GetURL(path string) (string, error)
}

func NewStorageProvider(cfg *config.AppConfig) (StorageProvider, error) {
    switch cfg.StorageProvider {
    // case "s3":
        // return NewS3Storage(cfg) 
    default:
        return NewLocalStorage(cfg) 
    }
}