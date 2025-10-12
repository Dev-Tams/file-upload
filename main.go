package main

import (
	"fmt"
	"io"
	"os"

	"github.com/dev-tams/file-upload/auth"
	"github.com/dev-tams/file-upload/config"
	"github.com/dev-tams/file-upload/handlers"
	"github.com/dev-tams/file-upload/handlers/admin"
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
		files.GET("/", handlers.GetAllFile)
		files.GET("/:id", handlers.GetFile)
		files.POST("/upload", handlers.PostFile)
		files.DELETE("/delete/:id", handlers.DeleteFile)


		adminRoutes := api.Group("/admin")
		adminRoutes.Use(auth.Middleware(), auth.AdminOnly())
		{
			adminRoutes.GET("/files", admin.GetAllFiles)
			adminRoutes.GET("/files/:id", admin.GetFile)
			adminRoutes.DELETE("/files/:id", admin.DeleteFile)

			adminRoutes.GET("/users", admin.FetchUsers)
			adminRoutes.GET("/users:id", admin.FetchUser)
			adminRoutes.DELETE("/users/:id", admin.DeleteUser)
		}

	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf(" server running on port %s", os.Getenv(port))
	router.Run(":" + port)

}
