package handlers

import (
	"net/http"

	"github.com/dev-tams/file-upload/internal/config"
	"github.com/dev-tams/file-upload/internal/models"
	"github.com/gin-gonic/gin"
)
// GetDeletedUsers godoc
// @Summary Get all deleted users
// @Description Returns all users where deleted_at is not null.
// @Tags Deleted Users
// @Produce json
// @Success 200 {object} map[string]interface{} "list of deleted users"
// @Failure 404 {object} map[string]interface{} "deleted users not found"
// @Security BearerAuth
// @Router /api/deleted/users [get]

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

// GetDeletedUser godoc
// @Summary Get a deleted user by email
// @Description Fetch a soft-deleted user using their email.
// @Tags Deleted Users
// @Produce json
// @Param email path string true "User Email"
// @Success 200 {object} map[string]interface{} "deleted user"
// @Failure 404 {object} map[string]interface{} "user not found"
// @Security BearerAuth
// @Router /api/deleted/users/{email} [get]

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

// RestoreUser godoc
// @Summary Restore a deleted user
// @Description Restores a soft-deleted user and sets deleted_at to null.
// @Tags Deleted Users
// @Produce json
// @Param email path string true "User Email"
// @Success 200 {object} map[string]interface{} "user restored"
// @Failure 404 {object} map[string]interface{} "user not found"
// @Failure 500 {object} map[string]interface{} "restore failed"
// @Security BearerAuth
// @Router /api/deleted/users/{email}/restore [put]

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

// FetchDeletedFiles godoc
// @Summary Get all deleted files
// @Description Returns all files where deleted_at is not null.
// @Tags Deleted Files
// @Produce json
// @Success 200 {object} map[string]interface{} "deleted files"
// @Failure 404 {object} map[string]interface{} "files not found"
// @Security BearerAuth
// @Router /api/deleted/files [get]

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
// FetchDeletedFile godoc
// @Summary Get deleted file by ID
// @Description Fetch a soft-deleted file and preload its user.
// @Tags Deleted Files
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} map[string]interface{} "deleted file"
// @Failure 404 {object} map[string]interface{} "file not found"
// @Security BearerAuth
// @Router /api/deleted/files/{id} [get]

func FetchDeletedFile(ctx *gin.Context){
	var file models.File
	id := ctx.Param("id")
	
	if err := config.DB.
	Unscoped().
	Where("id = ? AND deleted_at IS NOT NULL", id).Preload("User").First(&file).Error;
	err != nil{
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "deleted file not found",
			"err": err.Error(),
		})
		return
	}


	ctx.JSON(http.StatusOK, gin.H{"deleted_file": file})
}
// RestoreFile godoc
// @Summary Restore a deleted file
// @Description Restores a soft-deleted file and resets deleted_at to null.
// @Tags Deleted Files
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} map[string]interface{} "file restored"
// @Failure 404 {object} map[string]interface{} "file not found"
// @Failure 500 {object} map[string]interface{} "restore failed"
// @Security BearerAuth
// @Router /api/deleted/files/{id}/restore [put]

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
	Where("id = ? ", id).Model(&file).Update("deleted_at", nil).Error;
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