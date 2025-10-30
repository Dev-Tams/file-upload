package repositories

import (
	"github.com/dev-tams/file-upload/internal/models"
	"gorm.io/gorm"
)

type FileRepository struct {
	db *gorm.DB
}

var file models.File

func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Create(file *models.File) error {
	return r.db.Create(file).Error
}

func (r *FileRepository) GetById(id, userID string) (*models.File, error) {
		err := r.db.
		Where("id = ? AND user_id = ?", id, userID).
		Preload("User").
		First(&file).Error

	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (r *FileRepository) Delete(delete *models.File) error {
	if err := r.db.Delete(delete).Error; err != nil {
		return err
	}
	return nil
}

func (r *FileRepository) GET(userID string) error{
	return r.db.Where("user_id = ?", userID).Order("uploaded_at DESC").Preload("User").Error

}
