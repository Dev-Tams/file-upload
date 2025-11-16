package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dev-tams/file-upload/internal/config"
	"github.com/dev-tams/file-upload/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// PromoteToAdmin godoc
// @Summary Promote a user to admin
// @Description Change a user's role to admin by their ID.
// @Tags Admin Users
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]interface{} "user promoted"
// @Failure 404 {object} map[string]interface{} "user not found"
// @Failure 500 {object} map[string]interface{} "failed to update role"
// @Security BearerAuth
// @Router /api/admin/users/{user_id}/promote [put]

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

//SeedAdminFromENV godoc
// @Summary Seed super admin from environment variables
// @Description Creates a super admin user using SUPER_ADMIN_EMAIL and SUPER_ADMIN_PASSWORD from .env if not already present.
// @Tags Admin Setup
// @Success 200 {string} string "admin created or already exists"
// @Failure 500 {object} map[string]string "failed to create admin"

func SeedAdminFromENV() {
	email := os.Getenv("SUPER_ADMIN_EMAIL")
	password := os.Getenv("SUPER_ADMIN_PASSWORD")

	email = strings.TrimSpace(email)
	password = strings.TrimSpace(strings.Trim(password, `"'`))


	if email == "" || password == "" {
		fmt.Println("Missing SUPER_ADMIN_EMAIL or SUPER_ADMIN_PASSWORD in .env — skipping admin seeding")
		return
	}

	var user models.User
	id := uuid.New().String()

	err := config.DB.Where("email = ?", email).First(&user).Error
	if err == nil {
		fmt.Println("super admin already exists")
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("error checking admin:", err)
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
