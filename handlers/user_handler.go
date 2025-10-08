package handlers

import (
	"net/http"

	"github.com/dev-tams/file-upload/config"
	"github.com/dev-tams/file-upload/models"
	"github.com/gin-gonic/gin"
)


func FetchUsers(c *gin.Context) {

	var users []models.User

	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All users from DB",
		"data":    users,
	})

}

func FetchUser( c *gin.Context) {
	var user models.User

	id  := c.Param("id")
	if err := config.DB.Where(" id = ?", id).First(&user).Error; err != nil{
		c.JSON(http.StatusNotFound, gin.H{
			"err": err.Error(),
		})
		return
	}
	if err := config.DB.Find(&user).Error; err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"err" : err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user" : user,
	})

}
func DeleteUser(c *gin.Context) {

	var user models.User

	id := c.Param("id")

	if err := config.DB.Where("id = ? ", id).First(&user).Error; err != nil {
		c.JSONP(http.StatusNotAcceptable, gin.H{
			"error":   "user not found",
			"details": err.Error(),
		})
		return
	}
	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSONP(http.StatusNotAcceptable, gin.H{
			"error":   "failed to delete user",
			"details": err.Error(),
		})
	}

	c.JSONP(http.StatusNoContent, gin.H{})
}
