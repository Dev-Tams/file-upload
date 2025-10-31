package router

import (
	"github.com/dev-tams/file-upload/internal/handlers"
	admin "github.com/dev-tams/file-upload/internal/handlers/admin"
	auth "github.com/dev-tams/file-upload/internal/middleware"
	"github.com/dev-tams/file-upload/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, service *services.Service) {
	fileHandler := handlers.NewFileHandler(service)
	adminFHandler := admin.NewAdminFileHandler(service)
	adminUHandler := admin.NewUserHandler(service)

	api := r.Group("api")

	{
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)

		files := api.Group("/files")
		files.Use(auth.Middleware())
		{
			files.POST("/upload", fileHandler.PostFile)
			files.GET("/", fileHandler.GetAllFile)
			files.GET("/:id", fileHandler.GetFile)
			files.GET("/:id/download", fileHandler.DownloadFile)
			files.DELETE("/:id", fileHandler.DeleteFile)

		}

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
}
