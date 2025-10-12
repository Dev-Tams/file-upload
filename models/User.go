package models

import "time"

type User struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	Role     string `gorm:"default:user"`
	CreatedAt time.Time
}
