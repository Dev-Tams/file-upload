package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	Role      string `gorm:"default:user"`
	StorageLimit int  `gorm:"default:100"`
	Plan	string  `gorm:"default:free"`
	UsedStorage int 
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
