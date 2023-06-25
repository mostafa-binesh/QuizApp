package models

type Subject struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Title string `json:"title"`
	// relationships
	CourseID uint      `json:"course_id"`
	Course   *Course   `json:"course,omitempty"`
	Systems  []*System `json:"systems" gorm:"foreignKey:SubjectID"`
}
type SubjectWithSystems struct {
	ID      uint      `json:"id" gorm:"primary_key"`
	Title   string    `json:"title"`
	Systems []*System `json:"systems"`
}
