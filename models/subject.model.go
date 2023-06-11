package models

type Subject struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Title string `json:"title"`
}
