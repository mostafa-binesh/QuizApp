package models

type Question struct {
	ID        uint     `json:"id" gorm:"primary_key"`
	CourseID  uint     `json:"courseID"`
	Course    User     `json:"course" gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE;OnDelete:CSCADE"`
	Title     string   `json:"title"`
	Options   []Option `json:"options"`
	SubjectID uint     `json:"QuestionID"`
	Subject   Subject  `json:"subject" gorm:"foreignKey:QuestionID;constraint:OnUpdate:CASCADE;OnDelete:CSCADE"`
}
