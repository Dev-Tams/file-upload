package handlers

import (
	"net/http"
	"os"

	// "github.com/dev-tams/file-upload/models"
	"github.com/gin-gonic/gin"
)

// oute → GET /files/:filename
// What it should do:

// Read the filename parameter from the URL.

// Check if the file exists in uploads/.

// If exists → return the file as response.

// If not → return JSON with an error message.
func GetFile(c *gin.Context) {
	filename := c.Param("filename")
	if filename == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error": "filename param required"})
		return
	}

	filePath := "uploads/" + filename

	if _, err := os.Stat(filePath); os.IsNotExist(err){
		c.JSON(http.StatusBadRequest, gin.H{"error": "file not found" })
		return
	}
	c.File(filePath)
}

// Scan the uploads/ folder.

// Collect file metadata (e.g., name, size, modified time).

// Return as a JSON array.
func GetAllFile( c *gin.Context){
	dir  := "uploads"

	files, err := os.ReadDir(dir)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "problem with folser"})
		return
	}

	var fileList []gin.H
	for _, f := range files {
		info, _ := f.Info()
		fileList = append(fileList, gin.H{
			"name" : f.Name(),
			"type" : f.Type(),
			"size" : info.Size(),
			"modtime" : info.ModTime(), 
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"files" : fileList,
	})
}

func PostFile(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil{
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.Param(file.Filename)
	
	err = c.SaveUploadedFile(file, "uploads/" + file.Filename)
	if err != nil{
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully!", "file": file.Filename})
}