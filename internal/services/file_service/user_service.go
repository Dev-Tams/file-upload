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

func (u *Service) GetUserStorage(userID string) (*utils.UserStorageDTO, error) {
	if err := u.repo.GetUserStorage(userID); err != nil {
		return nil, err
	}
	user, err := u.repo.GetUserById(userID)
	if err != nil {
		return nil, err
	}
	dto := utils.FromUserStorageModel(*user)
	return &dto, nil
}

func (u *Service) GetAllUserStorage() ([]utils.UserStorageDTO, error) {
	users, err := u.repo.GetAllUserStorage()
	if err != nil {
		return nil, err
	}

	dto := utils.FromUserStorageModels(users)
	return dto, nil

}

func (u *Service) UpdateUserStorage(userID string, newLimit int64) (*models.User, error) {
	// newLimitMB from input
	limitBytes := newLimit * 1024 * 1024
	if err := u.repo.UpdateUserStorage(userID, limitBytes); err != nil {
		return nil, err
	}
	user, err := u.repo.GetUserById(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
