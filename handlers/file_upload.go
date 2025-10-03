package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

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

	if err := config.DB.Where("id = ?", ID).First(&file).Error; err != nil{
		c.JSON(http.StatusNotFound, gin.H{"error": "error finding file in db"})
		return
	}

	filePath := filepath.Join("uploads", file.StoredName)

	if _, err:= os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": " file not found on disk"})
	}

	// Serve the file
	c.File(filePath)

}

// Scan the uploads/ folder.

// Collect file metadata (e.g., name, size, modified time).

// Return as a JSON array.
func GetAllFile(c *gin.Context) {
	dir := "uploads"

	files, err := os.ReadDir(dir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "problem with folser"})
		return
	}

	var fileList []gin.H
	for _, f := range files {
		info, _ := f.Info()
		fileList = append(fileList, gin.H{
			"name":    f.Name(),
			"type":    f.Type(),
			"size":    info.Size(),
			"modtime": info.ModTime(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"files": fileList,
	})
}

func PostFile(c *gin.Context) {
	var files models.File

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	id := uuid.New().String()
	storedName := id + filepath.Ext(file.Filename)

	savedPath := filepath.Join("uploads", storedName)
	if err = c.SaveUploadedFile(file, savedPath); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	files.ID = id
	files.StoredName = storedName
	files.OriginalName = file.Filename
	files.DisplayName = file.Filename
	files.UploadedAt = time.Now()
	files.Size = file.Size
	files.Path = savedPath

	if err := config.DB.Create(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error saving file to db": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully!",
		"file":    files,
	})

}

// func FormatFile(path string) (newPath string, err error) {
//     // 1. Open the file on disk
//     // 2. Apply formatting (rename, resize, compress…)
//     // 3. Save changes
//     // 4. Return new path (or same path if unchanged)
// }
