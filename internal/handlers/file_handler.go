package handlers

import (
	"net/http"

	"github.com/dev-tams/file-upload/internal/config"
	"github.com/dev-tams/file-upload/internal/models"
	"github.com/dev-tams/file-upload/internal/services"
	"github.com/dev-tams/file-upload/internal/utils"
	"github.com/gin-gonic/gin"
)

type FileHandler struct {
    service *services.FileService
}

func NewFileHandler(service *services.FileService) *FileHandler {
    return &FileHandler{service: service}
}

func(f *FileHandler) GetFile(c *gin.Context) {
	ID := c.Param("id")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	
	userID := c.GetString("user_id")

    file, err := f.service.GetFile(ID, userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
        return
    }

	c.JSON(http.StatusOK, gin.H{"file": file})

}

func(f *FileHandler) GetAllFile(c *gin.Context) {

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

func (f *FileHandler) PostFile(c *gin.Context) {
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
	

	uploadedFiles, err := f.service.UploadFiles(userID, multipleFiles)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error uploading files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Files uploaded successfully!",
		"files":   uploadedFiles,
	})
}

func(f *FileHandler) DownloadFile(c *gin.Context) {
	ID := c.Param("id")

	switch ID {
	case "":
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID required"})
		return
	}

	userID := c.GetString("user_id")
	

	file, err := f.service.DownloadFile(ID, userID); if err != nil{
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}


	c.FileAttachment(file.Path, file.OriginalName)
}

func(f *FileHandler) DeleteFile(c *gin.Context) {
	ID := c.Param("id")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	userID := c.GetString("user_id")

	err := f.service.DeleteFile(ID, userID)
	if err != nil{
		c.JSON(http.StatusNotAcceptable, gin.H{
			"error":   "failed to delete file",
			"details": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)

}