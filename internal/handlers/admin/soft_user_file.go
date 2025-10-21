package handlers

import (
	"net/http"

	"github.com/dev-tams/file-upload/internal/config"
	"github.com/dev-tams/file-upload/internal/models"
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
	email := ctx.Param("email")
	if err := config.DB.Unscoped().
	Where("email = ? AND deleted_at IS NOT NULL", email).First(&user).Error;
	err != nil{
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "deleted user not found",
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"deleted_user": user})
}

func RestoreUser(ctx *gin.Context){

	var user models.User
	email := ctx.Param("email")
	if err := config.DB.Unscoped().
	Where("email = ? AND deleted_at IS NOT NULL", email).First(&user).Error;
	err != nil{
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "deleted user not found",
			"err": err.Error(),
		})
		return
	}

	if err := config.DB.Unscoped().
	Model(&user).
	Update("deleted_at", nil).Error;
	err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err" : "failed to restore user",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user restored successfully",
		"user":    user,
	})
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
func RestoreFile(ctx *gin.Context){
	var file models.File
	id := ctx.Param("id")
	
	if err := config.DB.
	Unscoped().
	Where("id = ? AND deleted_at IS NOT NULL", id).
	Preload("User").
	First(&file).Error;


	err != nil{
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "deleted file not found",
			"err": err.Error(),
		})
		return
	}

	if err := config.DB.
	Unscoped().
	Where("id = ? ", id).Model(file).Update("deleted_at", nil).Error;
	err != nil{
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "failed to restore file",
			"err": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "file restored successfully",
		"file":    file,
	})
}