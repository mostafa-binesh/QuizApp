package models

type Option struct {
	ID uint `json:"id" gorm:"primary_key"`
	Title string `json:"option"`
	Index string `json:"index"`
	QuestionID uint   `json:"QuestionID"`
	Question   User   `json:"course" gorm:"foreignKey:QuestionID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	IsCorrect bool `json:"status"`
}
