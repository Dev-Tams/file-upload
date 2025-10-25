package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dev-tams/file-upload/internal/config"
	"github.com/dev-tams/file-upload/internal/models"
	"github.com/dev-tams/file-upload/internal/repositories"
	"github.com/dev-tams/file-upload/internal/services"
	"github.com/dev-tams/file-upload/internal/utils"
	"github.com/gin-gonic/gin"
)



func GetFile(c *gin.Context) {
	ID := c.Param("id")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	var file models.File

	userID := c.GetString("user_id")

	if err := config.DB.Where("id = ? AND user_id = ?", ID, userID).Preload("User").First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	filePath := filepath.Join("uploads", file.StoredName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": " file not found on disk"})
	}

	// Serve the file
	fileDto := utils.FromFileModel(file)
	c.JSON(http.StatusOK, gin.H{"file": fileDto})

}

func GetAllFile(c *gin.Context) {

	userID := c.GetString("user_id")
	db := config.DB.Where("user_id = ?", userID).Order("uploaded_at DESC").Preload("User")

	pagination, err := utils.Paginate(c, db, models.File{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "pagination error"})
		return
	}

	files, ok := pagination.Data.([]models.File)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "type assertion failed"})
		return
	}
	fileDto := utils.FromFileModels(files)
	pagination.Data = fileDto

	c.JSON(http.StatusOK, gin.H{
		"files": gin.H{
			"page":       pagination.Page,
			"limit":      pagination.Limit,
			"total":      pagination.Total,
			"totalPages": pagination.TotalPages,
			"nextPage":   pagination.NextPage,
			"prevPage":   pagination.PrevPage,
			"data":       fileDto,
		},
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

	userID := c.GetString("user_id")
	

	repo := repositories.NewFileRepository(config.DB)
	service := services.NewFileService(repo)

	uploadedFiles, err := service.UploadFiles(userID, multipleFiles)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error uploading files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Files uploaded successfully!",
		"files":   uploadedFiles,
	})
}

func DownloadFile(c *gin.Context) {
	ID := c.Param("id")

	switch ID {
	case "":
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID required"})
		return
	}

	userID := c.GetString("user_id")
	var file models.File

	if err := config.DB.Where("id = ? AND user_id = ?", ID, userID).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	safePath := filepath.Clean(filepath.Join("uploads", file.StoredName))
	if !strings.HasPrefix(safePath, "uploads") {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid file path"})
		return
	}
	filePath := filepath.Join("uploads", file.StoredName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found on disk"})
		return
	}

	c.FileAttachment(filePath, file.OriginalName)
}

func DeleteFile(c *gin.Context) {
	var file models.File

	ID := c.Param("id")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	userID := c.GetString("user_id")

	if err := config.DB.Where("id = ? AND user_id = ?", ID, userID).First(&file).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"details": err.Error(),
		})
		return
	}

	if err := config.DB.Delete(&file).Error; err != nil {
		c.JSONP(http.StatusNotAcceptable, gin.H{
			"error":   "failed to delete file",
			"details": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)

}

// func FormatFile(path string) (newPath string, err error) {
//     // 1. Open the file on disk
//     // 2. Apply formatting (rename, resize, compress…)
//     // 3. Save changes
//     // 4. Return new path (or same path if unchanged)
// }
