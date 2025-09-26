package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Read the file from the request (you’ll typically specify a form field name, e.g., "file").

// Save it with a chosen filename into your uploads/ folder.

// Send a JSON response confirming success (and maybe return the file path or name for reference).
func PostFile(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil{
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	err = c.SaveUploadedFile(file, "uploads/" + file.Filename)
	if err != nil{
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully!", "file": file.Filename})
}