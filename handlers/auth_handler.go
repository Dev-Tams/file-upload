package handlers

import (
	"net/http"

	"github.com/dev-tams/file-upload/actions"
	"github.com/dev-tams/file-upload/config"
	"github.com/dev-tams/file-upload/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON format",
			"details": err.Error(),
		})
		return
	}

	if !actions.IsValidEmail(user.Email)  {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		return
	}

	if len(user.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 6 characters"})
		return
	}

	id  := uuid.New().String()
	var existingUser models.User
	if err := config.DB.Select("id").Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "email already registered",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}
	user.Password = string(hash)

	user.ID = id
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "could not create user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created",
	})

}
