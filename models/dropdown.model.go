package models

type DropDown struct {
	ID      uint     `json:"no" gorm:"primary_key"`
	Options []Option `json:"options,omitempty"`
	QuestionID uint     `json:"-"`
	Question   *Question  `json:"question,omitempty" gorm:"foreignKey:QuestionID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
}
