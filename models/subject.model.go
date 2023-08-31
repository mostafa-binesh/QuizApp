package models

type Subject struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Title string `json:"title"`
	// relationships
	CourseID uint      `json:"courseID"`
	Course   *Course   `json:"course,omitempty"`
	Systems  []*System `json:"systems" gorm:"foreignKey:SubjectID"`
}
type SubjectWithSystems struct {
	ID       uint      `json:"id" gorm:"primary_key"`
	Title    string    `json:"title"`
	Systems  []*System `json:"systems"`
	CourseID uint      `json:"courseID"`
}
type SubjectWithQuestionsCount struct {
	ID             uint                        `json:"id" gorm:"primary_key"`
	Title          string                      `json:"title"`
	Systems        []*SystemWithQuestionsCount `json:"systems"`
	CourseID       uint                        `json:"courseID"`
	QuestionsCount int                         `json:"questionsCount"`
}
