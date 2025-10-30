package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/dev-tams/file-upload/internal/config"
	"github.com/dev-tams/file-upload/internal/handlers"
	admin "github.com/dev-tams/file-upload/internal/handlers/admin"
	auth "github.com/dev-tams/file-upload/internal/middleware"
	"github.com/dev-tams/file-upload/internal/repositories"
	"github.com/dev-tams/file-upload/internal/services"
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

	SeedAd := flag.Bool("seed-admin", false, "Run database seeding for the admin user.")
	flag.Parse()

	config.ConnectDatabase()

	if *SeedAd {
		admin.SeedAdminFromENV()
		log.Println("Admin user seeding complete. Exiting.")
		return
	}

	repo := repositories.NewDbRepository(config.DB)
	service := services.NewService(repo)
	fileHandler := handlers.NewFileHandler(service)
	adminFHandler := admin.NewAdminFileHandler(service)
	adminUHandler := admin.NewUserHandler(service)
	router := gin.Default()
	router.Use(cors.Default())

	{
		api := router.Group("api")
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)

		files := api.Group("/files")
		files.Use(auth.Middleware())
		files.POST("/upload", fileHandler.PostFile)
		files.GET("/", fileHandler.GetAllFile)
		files.GET("/:id", fileHandler.GetFile)
		files.GET("/:id/download", fileHandler.DownloadFile)
		files.DELETE("/:id", fileHandler.DeleteFile)

		adminRoutes := api.Group("/admin")
		adminRoutes.Use(auth.Middleware(), auth.AdminOnly())
		{
			adminRoutes.GET("/users", adminUHandler.FetchUsers)
			adminRoutes.GET("/users/:user_id", adminUHandler.FetchUser)
			adminRoutes.DELETE("/users/:user_id", adminUHandler.DeleteUser)

			adminRoutes.GET("/users/:user_id/files", adminFHandler.GetAllFiles)
			adminRoutes.GET("/users/:user_id/files/:id", adminFHandler.GetFile)
			adminRoutes.GET("/users/:user_id/files/:id/download", adminFHandler.DownloadFile)
			adminRoutes.DELETE("/users/:user_id/files/:id", adminFHandler.DeleteFile)

			adminRoutes.PUT("/users/:user_id", admin.PromoteToAdmin)

			adminRoutes.GET("/users/deleted", admin.GetDeletedUsers)
			adminRoutes.GET("/users/deleted/:id", admin.GetDeletedUser)
			adminRoutes.GET("/users/deleted/files", admin.FetchDeletedFiles)
			adminRoutes.GET("/users/deleted/files/:id", admin.FetchDeletedFile)

			adminRoutes.PATCH("/users/restore/:email", admin.RestoreUser)
			adminRoutes.PATCH("/users/restore/files/:id", admin.RestoreFile)

			adminRoutes.PUT("reset-pasword", admin.ResetPassword)

		}

	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf(" server running on port %s", port)
	router.Run(":" + port)

}
