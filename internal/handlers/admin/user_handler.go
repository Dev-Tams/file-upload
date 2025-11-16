package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dev-tams/file-upload/internal/config"
	"github.com/dev-tams/file-upload/internal/models"
	"github.com/dev-tams/file-upload/internal/services/file_service"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	service *file_service.Service
}

func NewUserHandler(service *file_service.Service) *UserHandler {
	return &UserHandler{service: service}
}

// FetchUsers godoc
// @Summary Get all users
// @Description Returns all user records from the database.
// @Tags Users
// @Produce json
// @Success 200 {object} map[string]interface{} "All users from DB"
// @Failure 500 {object} map[string]string "server error"
// @Security BearerAuth
// @Router /api/users [get]

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

// FetchUser godoc
// @Summary Get user by ID
// @Description Fetch a single user using their user_id.
// @Tags Users
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]interface{} "user data"
// @Failure 404 {object} map[string]string "user not found"
// @Security BearerAuth
// @Router /api/users/{user_id} [get]

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

// DeleteUser godoc
// @Summary Delete a user
// @Description Deletes a user by their ID.
// @Tags Users
// @Param id path string true "User ID"
// @Success 204 "no content"
// @Failure 500 {object} map[string]interface{} "failed deleting user"
// @Security BearerAuth
// @Router /api/users/{id} [delete]
func (u *UserHandler) DeleteUser(c *gin.Context) {

	userID := c.Param("id")

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
// GetUserStorage godoc
// @Summary Get user storage
// @Description Fetch a user's storage usage and limit.
// @Tags Storage
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]interface{} "storage info"
// @Failure 404 {object} map[string]string "user not found"
// @Security BearerAuth
// @Router /api/users/{user_id}/storage [get]

func (u *UserHandler) GetUserStorage(c *gin.Context) {

	userID := c.Param("user_id")
	user, err := u.service.GetUserStorage(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data" : "user storage",
		"storage": user,

	})
}
// GetAllUserStorage godoc
// @Summary Get storage info for all users
// @Description Returns storage usage and limits for every user.
// @Tags Storage
// @Produce json
// @Success 200 {object} map[string]interface{} "list of user storage"
// @Failure 404 {object} map[string]string "no users found"
// @Security BearerAuth
// @Router /api/users/storage [get]

func (u *UserHandler) GetAllUserStorage(c *gin.Context) {

	users, err := u.service.GetAllUserStorage()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "all users storage",
		"data":     users,
	})
}
// UpdateUserStorage godoc
// @Summary Update a user's storage limit
// @Description Admin can manually override the storage limit for a user.
// @Tags Storage
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param body body object{newLimit=int64} true "New storage limit in MB"
// @Success 200 {object} map[string]interface{} "updated storage limit"
// @Failure 400 {object} map[string]interface{} "invalid json"
// @Failure 404 {object} map[string]string "user not found"
// @Security BearerAuth
// @Router /api/users/{user_id}/storage [put]

func (u *UserHandler) UpdateUserStorage(c *gin.Context) {
	var req struct {
		NewLimit int64
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "invalid json format",
			"details": err.Error(),
		})
		return
	}

	userID := c.Param("user_id")
	user, err := u.service.UpdateUserStorage(userID, req.NewLimit)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user storage updated successfully",
		"email":   user.Email,
		"storage": fmt.Sprintf("%dMB", user.StorageLimit),
	})
}
// AssignPlan godoc
// @Summary Assign a subscription plan to a user
// @Description Updates a user's plan (free, pro, etc.) and adjusts storage limit.
// @Tags Users
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param body body object{plan=string} true "Plan name"
// @Success 200 {object} map[string]interface{} "plan updated"
// @Failure 400 {object} map[string]interface{} "invalid json"
// @Failure 404 {object} map[string]string "user not found"
// @Security BearerAuth
// @Router /api/users/{user_id}/plan [put]

func (u *UserHandler) AssignPlan(c *gin.Context) {

	var req struct {
		Plan string
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "invalid json format",
			"details": err.Error(),
		})
		return
	}
	userID := c.Param("user_id")

	 user, err := u.service.AssignPlan(userID, req.Plan)
	 if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user plan updated successfully",
		"email":   user.Email,
		"plan":  user.Plan,
		"storage": fmt.Sprintf("%dMB", user.StorageLimit),
	})
}

// basic for testing purposes.
// ResetPassword godoc
// @Summary Reset a user's password (testing only)
// @Description Updates the user's password directly — not recommended for production.
// @Tags Users
// @Accept json
// @Produce json
// @Param body body object{email=string,password=string} true "Reset password payload"
// @Success 200 {object} map[string]interface{} "password updated"
// @Failure 400 {object} map[string]interface{} "invalid json"
// @Failure 404 {object} map[string]string "user not found"
// @Failure 500 {object} map[string]interface{} "failed to update"
// @Security BearerAuth
// @Router /api/users/reset-password [post]

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
