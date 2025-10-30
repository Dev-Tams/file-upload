package services

import (
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
func (f *Service) DeleteUser( userID string) error {
	user, err := f.repo.GetUserById(userID)
	if err != nil {
		return err
	}
	err = f.repo.DeleteUser(user)
	if err != nil {
		return err
	}
	return nil
}
