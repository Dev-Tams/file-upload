package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dev-tams/file-upload/internal/config"
	"github.com/dev-tams/file-upload/internal/models"
	"github.com/dev-tams/file-upload/internal/services/file_service"
	"github.com/dev-tams/file-upload/internal/utils"
	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	service *file_service.Service
}

func NewFileHandler(service *file_service.Service) *FileHandler {
	return &FileHandler{service: service}
}

// GetFile godoc
// @Summary Get a single file
// @Description Fetch a file belonging to the authenticated user by file ID.
// @Tags Files
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} map[string]interface{} "file object"
// @Failure 400 {object} map[string]string "id required"
// @Failure 404 {object} map[string]string "file not found"
// @Security BearerAuth
// @Router /api/files/{id} [get]

func (f *FileHandler) GetFile(c *gin.Context) {
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

// GetAllFile godoc
// @Summary Get all files for the authenticated user
// @Description Returns paginated list of user files with metadata.
// @Tags Files
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{} "paginated files"
// @Failure 500 {object} map[string]string "pagination error"
// @Security BearerAuth
// @Router /api/files [get]

func (f *FileHandler) GetAllFile(c *gin.Context) {

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

// PostFile godoc
// @Summary Upload one or multiple files
// @Description Upload files with idempotency protection. Requires `Idempotency-Key` header.
// @Tags Files
// @Accept mpfd
// @Produce json
// @Param Idempotency-Key header string true "Idempotency Key"
// @Param file formData file true "One or more files"
// @Success 200 {object} map[string]interface{} "Files uploaded successfully"
// @Failure 400 {object} map[string]string "missing Idempotency Key / no files uploaded"
// @Failure 500 {object} map[string]string "error uploading files"
// @Security BearerAuth
// @Router /api/files/upload [post]

func (f *FileHandler) PostFile(c *gin.Context) {

	// req idem key
	ctx := context.Background()
	idemKey := c.GetHeader("Idempotency-Key")

	if idemKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing Idempotency Key",
		})
		return
	}

	//cache if only redis is availabale
	if config.Cache != nil {
		//check idem key on cache
		if cached, err := config.Cache.Get(ctx, idemKey).Result(); err == nil {
			var cachedResp map[string]any

			if err := json.Unmarshal([]byte(cached), &cachedResp); err == nil {
				c.JSON(http.StatusOK, cachedResp)
				return
			}
		}
	}
	//req file
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//check file req
	multipleFiles := form.File["file"]
	if len(multipleFiles) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no files uploaded"})
		return
	}

	userID := c.GetString("user_id")

	//uploads file
	uploadedFiles, err := f.service.UploadFiles(userID, multipleFiles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error uploading files",
			"err" : err.Error(),
		})
		return
	}

	resp := gin.H{
		"message": "Files uploaded successfully!",
		"files":   uploadedFiles,
	}
	if config.Cache != nil {
		respJSON, _ := json.Marshal(resp)
		config.Cache.Set(ctx, idemKey, respJSON, 24*time.Hour)
	}
	c.JSON(http.StatusOK, resp)
}

// DownloadFile godoc
// @Summary Download a file
// @Description Downloads a file owned by the authenticated user.
// @Tags Files
// @Produce application/octet-stream
// @Param id path string true "File ID"
// @Success 200 "File attachment"
// @Failure 400 {object} map[string]string "ID required"
// @Failure 404 {object} map[string]string "file not found"
// @Security BearerAuth
// @Router /api/files/download/{id} [get]

func (f *FileHandler) DownloadFile(c *gin.Context) {
	ID := c.Param("id")

	switch ID {
	case "":
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID required"})
		return
	}

	userID := c.GetString("user_id")

	file, err := f.service.DownloadFile(ID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	c.FileAttachment(file.Path, file.OriginalName)
}

// DeleteFile godoc
// @Summary Delete a file
// @Description Deletes a file belonging to the authenticated user.
// @Tags Files
// @Param id path string true "File ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "id required"
// @Failure 406 {object} map[string]interface{} "failed to delete file"
// @Security BearerAuth
// @Router /api/files/{id} [delete]

func (f *FileHandler) DeleteFile(c *gin.Context) {
	ID := c.Param("id")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	userID := c.GetString("user_id")

	err := f.service.DeleteFile(ID, userID)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"error":   "failed to delete file",
			"details": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)

}
