package models

type Question struct {
	ID uint `json:"no" gorm:"primary_key"`
	// CourseID uint    `json:"-"`
	// Course   *Course `json:"course,omitempty" gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	Title       string `json:"question"`
	Status      string `json:"-"`
	Description string `json:"description"`
	// relationships
	Options  []*Option `json:"options,omitempty"`
	SystemID uint      `json:"-"`
	System   *System   `json:"system,omitempty" gorm:"foreignKey:SystemID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
}

//
type QuestionList struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}
type AdminCreateQuestionInput struct {
	Question      string `json:"email" validate:"required"`
	Option1       string `json:"option1" validate:"required"`
	Option2       string `json:"option2" validate:"required"`
	Option3       string `json:"option3" validate:"required"`
	Option4       string `json:"option4" validate:"required"`
	CorrectOption uint   `json:"correct" validate:"required"`
	Description   string `json:"description" validate:"required"`
	SystemID      uint   `json:"systemID" validate:"required"`
}
