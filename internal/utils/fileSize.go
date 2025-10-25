package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
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
	if err != nil{
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