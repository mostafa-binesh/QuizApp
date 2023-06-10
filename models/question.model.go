package models

type Question struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	CourseID uint   `json:"courseID"`
	Course   User   `json:"course" gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	// relationships
	Options   []Option `json:"options"`
	SubjectID uint     `json:"QuestionID"`
	Subject   Subject  `json:"subject" gorm:"foreignKey:SubjectID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
}
// 
type QuestionList struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

