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
		if err := utils.ValidateFile(file, 1, []string{".png", ".jpg", ".jpeg"}); err != nil {
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

		if err := s.repo.Create(&fileModel); err != nil {
			return nil, err
		}

		dto := utils.FromFileModel(fileModel)
		uploadedFiles = append(uploadedFiles, dto)
	}
	return uploadedFiles, nil
}
