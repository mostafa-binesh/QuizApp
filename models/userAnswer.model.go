package models

type UserAnswer struct {
	ID         uint     `json:"id" gorm:"primary_key"`
	QuestionID uint     `json:"QuestionID"`
	Question   User     `json:"course" gorm:"foreignKey:QuestionID;constraint:OnUpdate:CASCADE;OnDelete:CSCADE"`
	Note       string   `json:"note"`
	IsMarked   bool     `json:"isMarked" gorm:"default:false;"`
	Submitted  bool     `json:"submitted"`
	UserID     uint     `json:"userID"`
	User       User     `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE;OnDelete:CSCADE"`
	QuizID     uint     `json:"quizID"`
	Quiz       User     `json:"quiz" gorm:"foreignKey:QuizID;constraint:OnUpdate:CASCADE;OnDelete:CSCADE"`
	Answers    []Answer `json:"answers"`
}
