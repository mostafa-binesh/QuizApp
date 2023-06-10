package models

type Answer struct {
	ID uint `json:"id" gorm:"primary_key"`
	UserAnswerID uint `json:"userAnsderId"`
	UserAnswers   UserAnswer `json:"user" gorm:"foreignKey:UserAnswerID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	OptionID uint `json:"optionId"`
	Option   Option `json:"option" gorm:"foreignKey:OptionID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
}
