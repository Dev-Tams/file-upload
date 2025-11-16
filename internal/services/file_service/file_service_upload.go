package file_service

import (
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/dev-tams/file-upload/internal/models"
	"github.com/dev-tams/file-upload/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (f *Service) UploadFiles(userID string, files []*multipart.FileHeader) ([]utils.FileResponseDTO, error) {
	var uploadedFiles []utils.FileResponseDTO

	err := f.repo.WithUserLock(userID, func(tx *gorm.DB, user *models.User) error {
		//get user storage
		currentUsage, err := f.repo.GetUserStorageUsage(userID)
		if err != nil {
			return nil
		}
		//calc storage of upload
		var totalNewSize int64
		for _, file := range files {
			totalNewSize += file.Size
		}

		//check storage + current file upload
		limitBytes := int64(user.StorageLimit) * 1024 * 1024

		if err := utils.EnsureWithinQuota(currentUsage, totalNewSize, limitBytes); err != nil {
			return err
		}


		//range files
		for _, file := range files {

			//validate file by ext and max limit upload
			if err := utils.ValidateFile(file); err != nil {
				return err
			}

			//set file stored name
			id := uuid.New().String()
			storedName := id + filepath.Ext(file.Filename)

			//save file in folder
			savedPath, err := f.storage.Save(file, storedName)
			if err != nil {
				return err
			}

			//set file details by model
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

			//create file in db
			if err := f.repo.Create(tx, &fileModel); err != nil {
				return err
			}

			//format file response
			dto := utils.FromFileModel(fileModel)
			uploadedFiles = append(uploadedFiles, dto)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return uploadedFiles, nil
}
