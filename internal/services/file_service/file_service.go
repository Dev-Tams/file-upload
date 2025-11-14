package file_service

import (
	"github.com/dev-tams/file-upload/internal/models"
	"github.com/dev-tams/file-upload/internal/repositories"
	"github.com/dev-tams/file-upload/internal/storage"
	"github.com/dev-tams/file-upload/internal/utils"

	"gorm.io/gorm"
)

type Service struct {
	repo    *repositories.DbRepository
	storage storage.StorageProvider
}

func NewService(repo *repositories.DbRepository, storage storage.StorageProvider) *Service {
	return &Service{repo: repo, storage: storage}
}

func (f *Service) GetAllFiles(id, userID string) ([]utils.FileResponseDTO, error) {
	files, err := f.repo.GetAllFiles(userID)
	if err != nil {
		return nil, err
	}
	dto := utils.FromFileModels(files)
	return dto, nil

}

func (f *Service) GetFile(id, userID string) (*utils.FileResponseDTO, error) {
	file, err := f.repo.GetById(id, userID)
	if err != nil {
		return nil, err
	}

	if _, err := utils.FindFilePath(file); err != nil {
		return nil, err
	}

	dto := utils.FromFileModel(*file)
	return &dto, nil
}

func (f *Service) DeleteFile(id, userID string) error {

	return f.repo.WithUserLock(userID, func(tx *gorm.DB, user *models.User) error {
		var file models.File
		if err := tx.Where("id = ? AND user_id = ?", id, userID).First(&file).Error; err != nil {
			return err
		}

		if err := tx.Delete(&file).Error; err != nil {
			return err
		}

		return nil
	})
}

func (f *Service) DownloadFile(id, userID string) (*models.File, error) {
	file, err := f.repo.GetById(id, userID)
	if err != nil {
		return nil, err
	}

	if _, err := utils.FindFilePath(file); err != nil {
		return nil, err
	}

	return file, nil
}
