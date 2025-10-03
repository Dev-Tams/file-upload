package main

import (
	"fmt"
	"io"
	"os"

	"github.com/dev-tams/file-upload/config"
	"github.com/dev-tams/file-upload/handlers"
	"github.com/gin-gonic/gin"
)

func main() {

	gin.DisableConsoleColor()
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	config.ConnectDatabase()

	router := gin.Default()
	router.POST("/upload", handlers.PostFile)
	router.GET("/files/:id", handlers.GetFile)
	router.GET("/files", handlers.GetAllFile)

	fmt.Println(" server running on port 80")
	router.Run(":8080")
}
