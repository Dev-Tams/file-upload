package repositories

import (
	"github.com/dev-tams/file-upload/internal/models"
	"gorm.io/gorm"
)

type DbRepository struct {
	db *gorm.DB
}

var file models.File

func NewDbRepository(db *gorm.DB) *DbRepository {
	return &DbRepository{db: db}
}

func (r *DbRepository) Create(file *models.File) error {
	return r.db.Create(file).Error
}

func (r *DbRepository) GetById(id, userID string) (*models.File, error) {
		err := r.db.
		Where("id = ? AND user_id = ?", id, userID).
		Preload("User").
		First(&file).Error

	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (r *DbRepository) Delete(delete *models.File) error {
	if err := r.db.Delete(delete).Error; err != nil {
		return err
	}
	return nil
}

func (r *DbRepository) GetAllFiles(userID string) ([]models.File, error) {
	var files []models.File
	err := r.db.Where("user_id = ?", userID).Order("uploaded_at DESC").Preload("User").Find(&files).Error
	if err != nil {
		return nil, err
	}
	return files, nil
}
