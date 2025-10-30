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

type FileService struct {
	repo *repositories.FileRepository
}

func NewFileService(repo *repositories.FileRepository) *FileService {
	return &FileService{repo: repo}
}

func (f *FileService) UploadFiles(userID string, files []*multipart.FileHeader) ([]utils.FileResponseDTO, error) {

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

func (f *FileService) GetAllFiles(id, userID string) error {
     err := f.repo.GET(userID)
    if err != nil {
       return err
    }
	return nil

}
func (f *FileService) GetFile(id, userID string) (*utils.FileResponseDTO, error) {
    file, err := f.repo.GetById(id, userID)
    if err != nil {
        return nil, err
    }

    dto := utils.FromFileModel(*file)
    return &dto, nil
}
func (f *FileService) DeleteFile(id, userID string) error {
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

func (f *FileService) DownloadFile(id, userID string) (*models.File, error) {
	file, err := f.repo.GetById(id, userID)
	if err != nil {
		return nil, err
	}

	if _, err := utils.FindFilePath(file); err != nil {
		return nil, err
	}

	return file, nil
}
