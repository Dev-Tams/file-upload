package repositories

import (
	"github.com/dev-tams/file-upload/internal/models"
	"gorm.io/gorm"
)

type FileRepository struct {
	db *gorm.DB
}

func NewFileRepository (db *gorm.DB) *FileRepository{
	return &FileRepository{db: db}
}

func (r *FileRepository) Create(file *models.File) error{
	return r.db.Create(file).Error
}


func (r *FileRepository) GetById(userID, id string, find any) (any, error) {

	err :=  r.db.Where("id = ? AND user_id = ?", id, userID).Preload("User").First(&find).Error
	if err != nil{
		return nil, err
	}
	return &find, nil

}