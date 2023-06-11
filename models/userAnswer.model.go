package models

type UserAnswer struct {
	ID         uint      `json:"id,omitempty" gorm:"primary_key"`
	QuestionID uint      `json:"-"`
	Question   *Question `json:"question,omitempty" gorm:"foreignKey:QuestionID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	Note       *string   `json:"note,omitempty"`
	IsMarked   bool      `json:"isMarked,omitempty" gorm:"default:false;"`
	Submitted  bool      `json:"submitted,omitempty" gorm:"default:false;"`
	Status     string    `json:"status,omitempty"`
	SpentTime  uint      `json:"spentTime,omitempty"`

	UserID uint  `json:"-"`
	User   *User `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	QuizID uint  `json:"-"`
	Quiz   *Quiz `json:"quiz,omitempty" gorm:"foreignKey:QuizID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	// ! multiple-choice option
	// Answers    []*Answer `json:"answers,omitempty"`
	// ! single-choice option
	Answer *string `json:"answer,omitempty"`
}
