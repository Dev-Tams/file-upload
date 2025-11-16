package utils

import (
	"time"

	"github.com/dev-tams/file-upload/internal/models"
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


type UserStorageDTO struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	UsedStorage  int64  `json:"used_storage"`
	StorageLimit int64    `json:"storage_limit"`
	Plan         string `json:"plan"`
	Role         string `json:"role"`
	PercentUsed  float64 `json:"percent_used"`
}

func FromFileModel(file models.File) FileResponseDTO {
	var user models.User
	return FileResponseDTO{
		ID:          file.ID,
		DisplayName: file.DisplayName,
		Path:        file.Path,
		UploadedAt:  file.UploadedAt,
		Size:        file.Size,
		UserID:      file.UserID,
		UserEmail:   user.Email,
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


func FromUserStorageModel(user models.User) UserStorageDTO {
	percent := 0.0
	if user.StorageLimit > 0 {
		percent = (float64(user.UsedStorage) / float64(user.StorageLimit)) * 100
	}

	return UserStorageDTO{
		ID:           user.ID,
		Email:        user.Email,
		UsedStorage:  user.UsedStorage,
		StorageLimit: user.StorageLimit,
		PercentUsed:  percent,
		Plan:         user.Plan,
		Role:         user.Role,
	}
}

func FromUserStorageModels(users []models.User) []UserStorageDTO {
	var dtos []UserStorageDTO
	for _, u := range users {
		dtos = append(dtos, FromUserStorageModel(u))
	}
	return dtos
}
