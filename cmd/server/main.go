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
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnv()
}
func main() {

	config.InitRedis()
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
	repo := repositories.NewDbRepository(config.DB)
	service := file_service.NewService(repo)

	r := gin.Default()
	r.Use(cors.Default())
	router.RegisterRoutes(r, service)

	port := config.Config.Port

	fmt.Printf(" server running on port %s", port)
	r.Run(":" + port)

}
