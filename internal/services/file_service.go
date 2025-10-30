package services

import (
	"mime/multipart"

	"path/filepath"
	"time"

	"github.com/dev-tams/file-upload/internal/models"
	"github.com/dev-tams/file-upload/internal/repositories"
	"github.com/dev-tams/file-upload/internal/utils"
	"github.com/google/uuid"
)

type Service struct {
	repo *repositories.DbRepository
}

func NewService(repo *repositories.DbRepository) *Service {
	return &Service{repo: repo}
}

func (f *Service) UploadFiles(userID string, files []*multipart.FileHeader) ([]utils.FileResponseDTO, error) {

	var uploadedFiles []utils.FileResponseDTO

	for _, file := range files {
		id := uuid.New().String()
		storedName := id + filepath.Ext(file.Filename)
		if err := utils.ValidateFile(file); err != nil {
			return nil, err
		}

		savedPath, err := utils.SaveUploadedFile(file, storedName)
		if err != nil {
			return nil, err
		}

		fileModel := models.File{
			ID:           id,
			StoredName:   storedName,
			OriginalName: file.Filename,
			DisplayName:  file.Filename,
			UploadedAt:   time.Now(),
			Size:         file.Size,
			Path:         savedPath,
			UserID:       userID,
		}

		if err := f.repo.Create(&fileModel); err != nil {
			return nil, err
		}

		dto := utils.FromFileModel(fileModel)
		uploadedFiles = append(uploadedFiles, dto)
	}
	return uploadedFiles, nil
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
    file, err := f.repo.GetById(id, userID)
    if err != nil {
        return  err
    }
    err = f.repo.Delete(file)
	if err != nil {
        return err
    }
	return nil
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
