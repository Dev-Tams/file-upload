package handlers

import (
	"log"
	"net/http"

	"github.com/dev-tams/file-upload/internal/config"
	"github.com/dev-tams/file-upload/internal/models"
	"github.com/dev-tams/file-upload/internal/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	service *services.Service
}

func NewUserHandler(service *services.Service) *UserHandler {
	return &UserHandler{service: service}
}

func (u *UserHandler) FetchUsers(c *gin.Context) {
	users, err := u.service.GetAllUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All users from DB",
		"data":    users,
	})

}

func (u *UserHandler) FetchUser(c *gin.Context) {
	userID := c.Param("user_id")
	user, err := u.service.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})

}
func (u *UserHandler) DeleteUser(c *gin.Context) {

	userID := c.Param("user_id")

	err := u.service.DeleteUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed deleting user",
			"details": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
//basic for testing purposes.
func ResetPassword(c *gin.Context) {
	var req struct {
		Email    string
		Password string
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "invalid json format",
			"details": err.Error(),
		})
		return
	}

	var user models.User

	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}
	user.Password = string(hash)
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to update password",
			"details": err.Error(),
		})
		return
	}

	log.Println("Password reset for", user.Email)
	c.JSON(http.StatusOK, gin.H{
		"message": "user password updated successfully",
		"user":    user.Email,
	})

}
