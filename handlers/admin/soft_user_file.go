package admin

import (
	"net/http"

	"github.com/dev-tams/file-upload/config"
	"github.com/dev-tams/file-upload/models"
	"github.com/gin-gonic/gin"
)

func GetDeletedUsers(ctx *gin.Context) {
	var users []models.User
	if err:= config.DB.Unscoped().Where("deleted_at IS NOT NULL").Find(&users).Error;
	err != nil{
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "deleted file not found",
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"deleted_users": users})
}

func GetDeletedUser(ctx *gin.Context){
	var user models.User
	id := ctx.Param("id")
	if err := config.DB.Unscoped().
	Where("id = ? AND deleted_at IS NOT NULL", id).First(&user).Error;
	err != nil{
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "deleted user not found",
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"deleted_user": user})
}


func FetchDeletedFiles(ctx *gin.Context) {
	var files []models.File
	if err := config.DB.Unscoped().Where("deleted_at IS NOT NULL").Find(&files).Error;
	err != nil{
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "deleted files not found",
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"deleted_files": files})
}

func FetchDeletedFile(ctx *gin.Context){
	var file models.File
	id := ctx.Param("id")
	if err := config.DB.
	Unscoped().
	Where("id = ? AND deleted_at IS NOT NULL", id).Preload("User").Error;
	err != nil{
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "deleted file not found",
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"deleted_file": file})
}