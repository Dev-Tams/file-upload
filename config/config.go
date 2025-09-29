package config

import (
	"fmt"

	"github.com/dev-tams/file-upload/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


var DB *gorm.DB
func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})

	if err != nil{
		fmt.Println("error connecting to db", err)
	}
	database.AutoMigrate(&models.File{})

	DB = database
}