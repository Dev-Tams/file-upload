package repositories

import (
	"github.com/dev-tams/file-upload/internal/models"
)

var user models.User

func (r *DbRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *DbRepository) GetUserById(userID string) (*models.User, error) {
	err := r.db.
		Where("id = ?", userID).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *DbRepository) DeleteUser(delete *models.User) error {
	if err := r.db.Delete(delete).Error; err != nil {
		return err
	}
	return nil
}
func (r *DbRepository) UpdateUserStorage(userID string, newLimit int) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("storage_limit", newLimit).Error

}

func (r *DbRepository) AssignPlan(userID, newPlan string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("plan", newPlan).Error
}

func (r *DbRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
