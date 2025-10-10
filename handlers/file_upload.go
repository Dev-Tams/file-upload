package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	dto "github.com/dev-tams/file-upload/DTO"
	"github.com/dev-tams/file-upload/actions"
	"github.com/dev-tams/file-upload/config"
	"github.com/dev-tams/file-upload/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetFile(c *gin.Context) {
	ID := c.Param("id")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	var file models.File

	userID := c.GetString("user_id")

	if err := config.DB.Where("id = ? AND user_id = ?", ID, userID).Preload("user").First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	filePath := filepath.Join("uploads", file.StoredName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": " file not found on disk"})
	}

	// Serve the file
	fileDto := dto.FromFileModel(file)
	c.JSON(http.StatusOK, gin.H{"file": fileDto})

}

func GetAllFile(c *gin.Context) {

	userID := c.GetString("user_id")
	var files []models.File
	db := config.DB.Where("user_id= ?", userID).Find(&files)
	pagination, err := actions.Paginate(c, db, files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files": pagination,
	})
}

func PostFile(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	multipleFiles := form.File["file"]
	if len(multipleFiles) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no files uploaded"})
		return
	}

	var uploadedFiles []dto.FileResponseDTO

	for _, file := range multipleFiles {
		id := uuid.New().String()
		storedName := id + filepath.Ext(file.Filename)

		// validate file size & ext
		if err := actions.ValidateFile(file, 1, []string{".png", ".jpg", ".jpeg"}); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		savedPath := filepath.Join("uploads", storedName)
		if err := c.SaveUploadedFile(file, savedPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		userID := c.GetString("user_id")
		uploadedFile := models.File{
			ID:           id,
			StoredName:   storedName,
			OriginalName: file.Filename,
			DisplayName:  file.Filename,
			UploadedAt:   time.Now(),
			Size:         file.Size,
			Path:         savedPath,
			UserID:       userID,
		}

		if err := config.DB.Create(&uploadedFile).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error saving file to db": err.Error()})
			return
		}

		if err := config.DB.Preload("User").First(&uploadedFile, "id = ?", uploadedFile.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load user info"})
			return
		}
		fileDto := dto.FromFileModel(uploadedFile)
		uploadedFiles = append(uploadedFiles, fileDto)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Files uploaded successfully!",
		"files":   uploadedFiles,
	})
}

// func FormatFile(path string) (newPath string, err error) {
//     // 1. Open the file on disk
//     // 2. Apply formatting (rename, resize, compress…)
//     // 3. Save changes
//     // 4. Return new path (or same path if unchanged)
// }
