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


func (r *FileRepository) GetById(id string) (*models.File, error) {

	var file models.File
	err :=  r.db.Preload("User").First(&file, "id = ?", id).Error
	if err != nil{
		return nil, err
	}
	return &file, nil

}