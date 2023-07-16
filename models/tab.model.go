package models

type Tab struct {
	ID         uint      `json:"no" gorm:"primary_key"`
	// Tables     []Table   `json:"tables"`
	QuestionID uint      `json:"-"`
	Question   *Question `json:"question,omitempty" gorm:"foreignKey:QuestionID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
}
