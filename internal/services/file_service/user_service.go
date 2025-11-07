package file_service

import (
	"github.com/dev-tams/file-upload/internal/models"
	"github.com/dev-tams/file-upload/internal/utils"
)

func (u *Service) GetAllUser() ([]utils.UserResponseDTO, error) {
	users, err := u.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	dto := utils.FromUserModels(users)
	return dto, nil

}
func (u *Service) GetUser(userID string) (*utils.UserResponseDTO, error) {
	user, err := u.repo.GetUserById(userID)
	if err != nil {
		return nil, err
	}

	dto := utils.FromUserModel(*user)
	return &dto, nil
}
func (u *Service) DeleteUser(userID string) error {
	user, err := u.repo.GetUserById(userID)
	if err != nil {
		return err
	}
	err = u.repo.DeleteUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *Service) UpdateUserStorage(userID string, newLimit int) (*models.User, error) {
	if err := u.repo.UpdateUserStorage(userID, newLimit); err != nil {
		return nil, err
	}
	user, err := u.repo.GetUserById(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
