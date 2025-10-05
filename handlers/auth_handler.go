package handlers

import (
	"net/http"

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
		})
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

func DeleteUser( c *gin.Context){

	var user models.User

	id := c.Param("id")
	
	if err := config.DB.Where("id = ? ", id).First(&user).Error; err != nil{
		c.JSONP(http.StatusNotAcceptable, gin.H{
			"error": "user not found",
			"details" : err.Error(),
		})
	}
	if err := config.DB.Delete(&user, id).Error; err != nil{
		c.JSONP(http.StatusNotAcceptable, gin.H{
			"error": "failed to delete user",
			"details" : err.Error(),
		})
	}
}
func FindUsers(c *gin.Context) {
	
	var users []models.User

	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All users from DB",
		"data":    users,
	})

}


