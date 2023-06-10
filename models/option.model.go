package models

type Option struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	Title     string `json:"option"`
	Index     string `json:"index"`
	IsCorrect bool   `json:"status"`
	// relationships
	QuestionID uint     `json:"QuestionID"`
	Question   Question `json:"course" gorm:"foreignKey:QuestionID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
}
