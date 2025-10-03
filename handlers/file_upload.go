package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

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
	// if err := config.DB.First(&files, ID); err != nil{
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "error finding file"})
	// 	return
	// }

	if err := config.DB.Where("id = ?", ID).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "error finding file in db"})
		return
	}

	filePath := filepath.Join("uploads", file.StoredName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": " file not found on disk"})
	}

	// Serve the file
	c.File(filePath)

}

// Scan the uploads/ folder.

// Collect file metadata (e.g., name, size, modified time).

// Return as a JSON array.
func GetAllFile(c *gin.Context) {

	var files []models.File
	if err := config.DB.Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch files"})
		return
	}

	var fileList []gin.H
	for _, f := range files {
		fileList = append(fileList, gin.H{
			"id":           f.ID,
			"originalName": f.OriginalName,
			"size":         f.Size,
			"displayName":  f.DisplayName,
			"path":         f.Path,
			"uploadedAt":   f.UploadedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"files": fileList,
	})
}

func PostFile(c *gin.Context) {
	form, err:= c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	multipleFiles := form.File["file"]
	if len(multipleFiles) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no files uploaded"})
		return
	}

	var uploadedFiles []models.File

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

		uploadedFile := models.File{
			ID:           id,
			StoredName:   storedName,
			OriginalName: file.Filename,
			DisplayName:  file.Filename,
			UploadedAt:   time.Now(),
			Size:         file.Size,
			Path:         savedPath,
		}

		if err := config.DB.Create(&uploadedFile).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error saving file to db": err.Error()})
			return
		}

		uploadedFiles = append(uploadedFiles, uploadedFile)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Files uploaded successfully!",
		"files":   uploadedFiles, // return all files, not just the last one
	})
}

// func FormatFile(path string) (newPath string, err error) {
//     // 1. Open the file on disk
//     // 2. Apply formatting (rename, resize, compress…)
//     // 3. Save changes
//     // 4. Return new path (or same path if unchanged)
// }
