package models

import (
	"time"
)

type File struct {
	ID          string         `gorm:"uuid"`
	StoredName  string         
	OriginalName string         
	DisplayName  string        
	UploadedAt   time.Time       
	Size        int64          
	Path        string 			`gorm:"not null;"`
	UserID     string 			`gorm:"not null;index"`
	User        User           `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}