package models

type File struct{
	ID int `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}