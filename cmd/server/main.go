// @title File Upload API
// @version 1.0
// @description API for file upload management
// @host localhost:8000
// @basePath /api
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Insert your JWT token like: Bearer <token>

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/dev-tams/file-upload/cmd/server/router"
	"github.com/dev-tams/file-upload/internal/config"
	admin "github.com/dev-tams/file-upload/internal/handlers/admin"
	"github.com/dev-tams/file-upload/internal/repositories"
	"github.com/dev-tams/file-upload/internal/services/file_service"
	"github.com/dev-tams/file-upload/internal/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	config.LoadEnv()
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

	config.InitRedis()
	cfg := config.Config
	storageProvider, err := storage.NewStorageProvider(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize storage provider: %v", err)
	}
	repo := repositories.NewDbRepository(config.DB)
	service := file_service.NewService(repo, storageProvider)

	r := gin.Default()
	r.Use(cors.Default())
	router.RegisterRoutes(r, service)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	port := config.Config.Port
	fmt.Printf(" server running on port %s", port)
	r.Run(":" + port)

}
