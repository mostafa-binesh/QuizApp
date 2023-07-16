package models

type Option struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	Title     string `json:"option,omitempty"`
	Index     string `json:"index,omitempty"`
	IsCorrect bool   `json:"status"`
	// relationships
	QuestionID uint      `json:"-"`
	Question   *Question `json:"question,omitempty" gorm:"foreignKey:QuestionID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
}
type AdminOptionInput struct {
	Title     string `form:"title" validate:"required"`
	IsCorrect bool   `form:"isCorrect" validate:"required"`
}
