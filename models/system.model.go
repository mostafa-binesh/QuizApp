package models

type System struct {
	ID        uint        `json:"id" gorm:"primary_key"`
	Title     string      `json:"title"`
	SubjectID uint        `json:"subjectID"`
	Subject   *Subject    `json:"subject,omitempty"`
	Questions []*Question `json:"questions,omitempty" gorm:"foreignKey:SystemID"`
}
type SystemWithQuestionsCount struct {
	ID             uint        `json:"id" gorm:"primary_key"`
	Title          string      `json:"title"`
	SubjectID      uint        `json:"subjectID"`
	QuestionsCount int        `json:"questionsCount"`
}
