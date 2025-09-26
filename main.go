package main

import (
	"github.com/gin-gonic/gin"
	"github.com/dev-tams/file-upload/handlers"
	// "net/http"

)

func main() {

	router := gin.Default()
	router.POST("/upload", handlers.PostFile)

	router.Run(":8080")
}