package utils

import (

	"github.com/dev-tams/file-upload/internal/models"
)

type UserResponseDTO struct {
	ID        string `json:"id"`
	UserEmail string `json:"user_email"`
	Plan	  string `json:"plan"`
	StorageLimit int64 `json:"storage_limit"` 
	UsedStorage	int64  `json:"used_storage"`
	Role      string `json:"role"`
}

func FromUserModel(user models.User) UserResponseDTO {
	return UserResponseDTO{
		ID:          user.ID,
		UserEmail:   user.Email,
		Plan:		 user.Plan,
		StorageLimit: user.StorageLimit,
		UsedStorage: user.UsedStorage,
		Role:        user.Role,	
		
	}
}

// for many users
func FromUserModels(users []models.User) []UserResponseDTO {
	var dtos []UserResponseDTO
	for _, u := range users {
		dtos = append(dtos, FromUserModel(u))
	}
	return dtos
}
