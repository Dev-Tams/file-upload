package config

import (
	"fmt"
	"log"
	"os"

	"github.com/dev-tams/file-upload/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB


func ConnectDatabase() {
		if err := os.MkdirAll("data", 0755); err != nil{
			fmt.Println("error creating folder", err.Error())
	}
	
	database, err := gorm.Open(sqlite.Open("data/file.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("error connecting to db", err)
	}
	database.AutoMigrate(&models.File{}, models.User{})

	DB = database
}
