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
	StorageLimit int64 `gorm:"default:104857600"` // 100 MB in bytes
	Plan	string  `gorm:"default:free"`
	UsedStorage int64 `gorm:"default:0"` 
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) LimitMB() float64 {
	return float64(u.StorageLimit) / (1024 * 1024)
}

func (u *User) UsedMB() float64 {
	return float64(u.UsedStorage) / (1024 * 1024)
}
