package config

import (
	"log"

	"github.com/dev-tams/file-upload/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


var DB *gorm.DB
func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
	if err != nil{
		log.Fatal("error connecting to db", err)
	}
	database.AutoMigrate(&models.File{})

	DB = database
}