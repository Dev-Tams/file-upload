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
		files.POST("/upload", handlers.PostFile)
		files.GET("/", handlers.GetAllFile)
		files.GET("/:id", handlers.GetFile)
		files.GET("/:id/download", handlers.DownloadFile)
		files.DELETE("/:id", handlers.DeleteFile)

		adminRoutes := api.Group("/admin")
		adminRoutes.Use(auth.Middleware(), auth.AdminOnly())
		{
			adminRoutes.GET("/users", admin.FetchUsers)
			adminRoutes.GET("/users/:user_id", admin.FetchUser)
			adminRoutes.DELETE("/users/:user_id", admin.DeleteUser)

			adminRoutes.GET("/users/:user_id/files", admin.GetAllFiles)
			adminRoutes.GET("/users/:user_id/files/:id", admin.GetFile)
			adminRoutes.GET("/users/:user_id/files/:id/download", admin.DownloadFile)
			adminRoutes.DELETE("/users/:user_id/files/:id", admin.DeleteFile)

			adminRoutes.PUT("/users/:user_id", admin.PromoteToAdmin)

		}

	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf(" server running on port %s", os.Getenv(port))
	router.Run(":" + port)

}
