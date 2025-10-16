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

func GetFile(ctx *gin.Context) {
	ID := ctx.Param("id")
	userID := ctx.Param("user_id")

	switch {
	case ID == "":
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	case userID == "":
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}
	var file models.File

	if err := config.DB.Where("user_id = ? AND id = ?", userID, ID).Preload("User").First(&file).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	filePath := filepath.Join("uploads", file.StoredName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": " file not found on disk"})
	}

	// Serve the file
	fileDto := dto.FromFileModel(file)
	ctx.JSON(http.StatusOK, gin.H{"file": fileDto})

}

func GetAllFiles(ctx *gin.Context) {

	userID := ctx.Param("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "user id required",
		})
	}

	var files []models.File
	db := config.DB.Where("user_id = ?", userID).
	Order("uploaded_at DESC").
	Preload("User")
	
	pagination, err := actions.Paginate(ctx, db, &models.File{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "pagination error"})
		return
	}
	
	files, ok := pagination.Data.([]models.File)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "type assertion failed"})
		return
	}
	fileDto := dto.FromFileModels(files)
	pagination.Data = fileDto


	ctx.JSON(http.StatusOK, gin.H{
		"files": gin.H{
		"page":        pagination.Page,
		"limit":       pagination.Limit,
		"total":       pagination.Total,
		"totalPages":  pagination.TotalPages,
		"nextPage":    pagination.NextPage,
		"prevPage":    pagination.PrevPage,
		"data":        fileDto,
	},
	})
}




func DeleteFile(ctx *gin.Context) {
	var file models.File

	ID := ctx.Param("id")
	if ID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}

	if err := config.DB.Preload("User").First(&file, "id = ?", ID).Error; err != nil {
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
