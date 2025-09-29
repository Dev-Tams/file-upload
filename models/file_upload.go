package models

 import "gorm.io/gorm"

type File struct {
	gorm.Model
	ID          string         `gorm:"uuid"`
	StoredName  string         
	OriginalName string         
	DisplayName  string        
	UploadedAt   string        
	Size        int64          
	Path        string
}