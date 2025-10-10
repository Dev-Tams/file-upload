package dto

import (
	"time"

	"github.com/dev-tams/file-upload/models"
)

type FileResponseDTO struct {
	ID          string    `json:"id"`
	DisplayName string    `json:"display_name"`
	Path        string    `json:"path"`
	UploadedAt  time.Time `json:"uploaded_at"`
	Size        int64     `json:"size"`
	UserID      string    `json:"user_id"`
	UserEmail   string    `json:"user_email"`
}

func FromFileModel(file models.File) FileResponseDTO {
	return FileResponseDTO{
		ID:          file.ID,
		DisplayName: file.DisplayName,
		Path:        file.Path,
		UploadedAt:  file.UploadedAt,
		Size:        file.Size,
		UserID:      file.UserID,
		UserEmail:   file.User.Email,
	}
}

// for many files
func FromFileModels(files []models.File) []FileResponseDTO {
	var dtos []FileResponseDTO
	for _, f := range files {
		dtos = append(dtos, FromFileModel(f))
	}
	return dtos
}
