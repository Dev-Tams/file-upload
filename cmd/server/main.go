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

	r := gin.Default()
	r.Use(cors.Default())
	router.RegisterRoutes(r, service)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf(" server running on port %s", port)
	r.Run(":" + port)

}
