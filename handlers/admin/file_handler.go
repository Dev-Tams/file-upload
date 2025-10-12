package admin

import (
	"net/http"
	"os"
	"path/filepath"

	dto "github.com/dev-tams/file-upload/DTO"
	"github.com/dev-tams/file-upload/actions"
	"github.com/dev-tams/file-upload/config"
	"github.com/dev-tams/file-upload/models"
	"github.com/gin-gonic/gin"
)

func GetFile(c *gin.Context) {
	ID := c.Param("id")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	var file models.File

	if err := config.DB.Preload("user").First(&file, "id = ?", ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
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

func GetAllFiles(ctx *gin.Context) {
	db := config.DB.Order("uploaded_at DESC")
	pagination, err := actions.Paginate(ctx, db, &models.File{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "pagination error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"files": pagination,
	})
}


func DeleteFile(ctx *gin.Context) {
	var file models.File

	ID := ctx.Param("id")
	if ID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}

	if err := config.DB.Preload("user").First(&file, "id = ?", ID).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"details": err.Error(),
		})
		return
	}

	if err := config.DB.Delete(&file).Error; err != nil {
		ctx.JSONP(http.StatusNotAcceptable, gin.H{
			"error":   "failed to delete file",
			"details": err.Error(),
		})
		return
	}

	ctx.Status(http.StatusNoContent)

}