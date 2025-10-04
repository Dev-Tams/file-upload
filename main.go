package main

import (
	"fmt"
	"io"
	"os"

	"github.com/dev-tams/file-upload/config"
	"github.com/dev-tams/file-upload/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

	func init() {
		godotenv.Load()
	}
func main() {

	gin.DisableConsoleColor()
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	config.ConnectDatabase()

	router := gin.Default()
	router.Use(cors.Default())

	{
		api := router.Group("api")
		api.POST("/upload", handlers.PostFile)
		api.GET("/files/:id", handlers.GetFile)
		api.GET("/files", handlers.GetAllFile)

	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}


	fmt.Println(" server running on port 8080")
	router.Run(":" + port)

}
