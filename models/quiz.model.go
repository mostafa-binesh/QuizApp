package models

type Quiz struct {
	ID     uint `json:"id" gorm:"primary_key"`
	UserID uint `json:"userID"`
	User   User `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE;OnDelete:CSCADE"`
	// TODO lesson >> ? lesson == course ?
	Status   string `json:"status"`
	CourseID uint   `json:"courseID"`
	Course   User   `json:"course" gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE;OnDelete:CSCADE"`
}
type QuizInput struct {
	QuestionsCount int `json:"questionsCount" validate:"required,min=1"`
}