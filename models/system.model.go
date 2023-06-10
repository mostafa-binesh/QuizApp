package models

type System struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	Title     string     `json:"title"`
	Questions []Question `json:"questions" gorm:"foreignKey:SystemID"`
	SubjectID uint       `json:"subjectID"`
	Subject   Subject    `json:"subject"`
}
