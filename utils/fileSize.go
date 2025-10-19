package actions

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
)

// Limit file size (c.Request.ContentLength).

// Restrict allowed types (check file.Header.Get("Content-Type") or file extension).
func ValidateFile(file *multipart.FileHeader, maxSize int, allowedExts[] string) error{

	limit := int64(maxSize) * 1024 * 1024
	if file.Size > limit {
		return fmt.Errorf(" file too large, max :%d MB",  maxSize)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	for _, e := range allowedExts {
        if ext == e {
            return nil // valid
        }
    }
	return fmt.Errorf("file type not allowed: %s", ext)
}
