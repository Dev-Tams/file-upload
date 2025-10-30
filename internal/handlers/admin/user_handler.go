package handlers

import (
	"net/http"

	"github.com/dev-tams/file-upload/internal/services"
	"github.com/gin-gonic/gin"

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
		"data":    user,
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
