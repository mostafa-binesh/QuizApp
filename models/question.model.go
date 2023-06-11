package models

type Question struct {
	ID       uint   `json:"no" gorm:"primary_key"`
	CourseID uint   `json:"-"`
	Course   *Course `json:"course,omitempty" gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	Title    string `json:"question"`
	Status   string `json:"-"`
	// relationships
	Options  []*Option `json:"options,omitempty"`
	SystemID uint     `json:"-"`
	System   *System   `json:"system,omitempty" gorm:"foreignKey:SystemID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
}

//
type QuestionList struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}
