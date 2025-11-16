package handlers

import (
	"net/http"

	"github.com/dev-tams/file-upload/internal/config"
	"github.com/dev-tams/file-upload/internal/models"
	"github.com/dev-tams/file-upload/internal/services/file_service"
	"github.com/dev-tams/file-upload/internal/utils"
	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	service *file_service.Service
}

func NewAdminFileHandler(service *file_service.Service) *FileHandler {
	return &FileHandler{service: service}
}
// GetFile godoc
// @Summary Get a file by ID
// @Description Fetch a single file belonging to a specific user.
// @Tags Admin Files
// @Produce json
// @Param user_id path string true "User ID"
// @Param id path string true "File ID"
// @Success 200 {object} map[string]interface{} "file object"
// @Failure 400 {object} map[string]interface{} "missing parameters"
// @Failure 404 {object} map[string]interface{} "file not found"
// @Security BearerAuth
// @Router /api/admin/users/{user_id}/files/{id} [get]

func(f *FileHandler) GetFile(ctx *gin.Context) {
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

	file, err := f.service.GetFile(ID, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"file": file})

}

// GetAllFiles godoc
// @Summary Get all files for a user (paginated)
// @Description Returns paginated list of files belonging to a user.
// @Tags Admin Files
// @Produce json
// @Param user_id path string true "User ID"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{} "paginated file data"
// @Failure 400 {object} map[string]interface{} "user id required"
// @Failure 500 {object} map[string]interface{} "pagination error"
// @Security BearerAuth
// @Router /api/admin/users/{user_id}/files [get]

func(f *FileHandler) GetAllFiles(ctx *gin.Context) {

	userID := ctx.Param("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": "user id required",
		})
		return
	}

	var files []models.File
	db := config.DB.Where("user_id = ?", userID).
		Order("uploaded_at DESC").
		Preload("User")

	pagination, err := utils.Paginate(ctx, db, models.File{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "pagination error"})
		return
	}

	files, ok := pagination.Data.([]models.File)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "type assertion failed"})
		return
	}
	fileDto := utils.FromFileModels(files)
	pagination.Data = fileDto

	ctx.JSON(http.StatusOK, gin.H{
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

// DownloadFile godoc
// @Summary Download a file
// @Description Download a file belonging to a user. Returns file attachment.
// @Tags Admin Files
// @Produce application/octet-stream
// @Param user_id path string true "User ID"
// @Param id path string true "File ID"
// @Success 200 {string} binary "file download"
// @Failure 400 {object} map[string]interface{} "missing parameters"
// @Failure 404 {object} map[string]interface{} "file not found"
// @Security BearerAuth
// @Router /api/admin/users/{user_id}/files/{id}/download [get]

func(f *FileHandler) DownloadFile(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.Param("user_id")

	switch {
	case id == "":
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	case userID == "":
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}

	
	file, err := f.service.DownloadFile(id, userID); if err != nil{
		ctx.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	ctx.FileAttachment(file.Path, file.OriginalName)
	ctx.Status(http.StatusOK)
}

// DeleteFile godoc
// @Summary Delete a file
// @Description Soft-deletes a file belonging to a specific user.
// @Tags Admin Files
// @Produce json
// @Param user_id path string true "User ID"
// @Param id path string true "File ID"
// @Success 204 {string} string "file deleted"
// @Failure 400 {object} map[string]interface{} "missing parameters"
// @Failure 406 {object} map[string]interface{} "delete failed"
// @Security BearerAuth
// @Router /api/admin/users/{user_id}/files/{id} [delete]

func(f *FileHandler) DeleteFile(ctx *gin.Context) {

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

	err := f.service.DeleteFile(ID, userID)
	if err != nil{
		ctx.JSON(http.StatusNotAcceptable, gin.H{
			"error":   "failed to delete file",
			"details": err.Error(),
		})
		return
	}

	ctx.Status(http.StatusNoContent)

}
