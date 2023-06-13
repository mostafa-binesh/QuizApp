package models

type System struct {
	ID        uint        `json:"id" gorm:"primary_key"`
	Title     string      `json:"title"`
	Questions []*Question `json:"questions,omitempty" gorm:"foreignKey:ID"`
	SubjectID uint        `json:"-"`
	Subject   *Subject    `json:"subject,omitempty"`
}
