package main

import (
	"fmt"
	"io"
	"os"

	"github.com/dev-tams/file-upload/auth"
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
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)

		files := api.Group("/files")
		files.Use(auth.Middleware())
		files.POST("/upload", handlers.PostFile)
		files.GET("/:id", handlers.GetFile)
		files.GET("/", handlers.GetAllFile)

		users := api.Group("/users")
		users.Use(auth.Middleware())
		users.GET("/", handlers.FetchUsers)
		users.GET("/:id", handlers.FetchUser)
		users.DELETE("/delete/:id", handlers.DeleteUser)

	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf(" server running on port %s", os.Getenv(port))
	router.Run(":" + port)

}
