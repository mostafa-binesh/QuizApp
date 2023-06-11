package models

type Answer struct {
	ID uint `json:"id" gorm:"primary_key"`
	UserAnswerID uint `json:"userAnsderId,omitempty"`
	UserAnswers   *UserAnswer `json:"user,omitempty" gorm:"foreignKey:UserAnswerID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	OptionID uint `json:"optionId,omitempty"`
	Option   *Option `json:"option,omitempty" gorm:"foreignKey:OptionID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
}
