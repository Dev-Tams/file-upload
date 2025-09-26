package main

import (
	"github.com/gin-gonic/gin"
	"github.com/dev-tams/file-upload/handlers"

)

func main() {

	router := gin.Default()
	router.POST("/upload", handlers.PostFile)
	router.GET("/files/:filename", handlers.GetFile)
	router.GET("/files", handlers.GetAllFile)

	router.Run(":8080")
}