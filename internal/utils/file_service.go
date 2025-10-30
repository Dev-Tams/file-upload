package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/dev-tams/file-upload/internal/models"
)

const (
	maxSize = 5 * 1024 * 1024
)

var allowedExts = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"application/pdf": true,
}

func ValidateFile(file *multipart.FileHeader) error {

	if file.Size > maxSize {
		return errors.New(" file too larage")
	}

	fileType := file.Header.Get("Content-Type")
	if !allowedExts[fileType] {
		return errors.New("unsupported file type")
	}
	return nil
}

func SaveUploadedFile(file *multipart.FileHeader, storedName string) (string, error) {

	err := os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		return "", fmt.Errorf(" error creating dir")
	}

	filePath := filepath.Join("uploads", storedName)

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf(" error opening file %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf(" error creating  file %w", err)
	}
	defer dst.Close()

	_, err = dst.ReadFrom(src)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return filePath, nil

}

func FindFilePath(file *models.File) (*models.File, error) {
    baseDir := "uploads"
    safePath := filepath.Join(baseDir, filepath.Clean(file.StoredName))

    rel, err := filepath.Rel(baseDir, safePath)
    if err != nil || strings.HasPrefix(rel, "..") {
        return nil, fmt.Errorf("invalid file path")
    }

    if _, err := os.Stat(safePath); os.IsNotExist(err) {
        return nil, fmt.Errorf("file not found")
    }

    return file, nil
}