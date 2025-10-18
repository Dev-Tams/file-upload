package admin

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dev-tams/file-upload/config"
	"github.com/dev-tams/file-upload/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func PromoteToAdmin(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	var user models.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user.Role = "admin"
	if err := config.DB.Save(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update role"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user promoted to admin"})
}

func SeedAdminFromENV() {
	email := os.Getenv("SUPER_ADMIN_EMAIL")
	password := os.Getenv("SUPER_ADMIN_PASSWORD")

	if email == "" || password == "" {
		fmt.Println("super admin config not set")
	}

	var user models.User
	id := uuid.New().String()

	if err := config.DB.Where("email =?", email).First(&user).Error; err != nil {
		fmt.Println(" super admin already exists")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("error hashing password")
		return
	}

	user = models.User{
		ID:       id,
		Email:    email,
		Password: string(hash),
		Role:     "admin",
	}

	if err = config.DB.Create(&user).Error; err != nil {
		fmt.Println("failed to create admin")
		return
	}
	log.Printf("admin created %s", email)
}
