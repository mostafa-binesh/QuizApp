package models

import (
	U "docker/utils"
	"gorm.io/gorm"
)

type Question struct {
	ID          uint     `json:"no" gorm:"primary_key"`
	Title       string   `json:"question"`
	Status      string   `json:"-"`
	Description string   `json:"description"`
	Image       []string `json:"image" gorm:"type:varchar(255)[]"`
	// relationships
	Options  []*Option `json:"options,omitempty"`
	SystemID uint      `json:"-"`
	System   *System   `json:"system,omitempty" gorm:"foreignKey:SystemID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	// although we could get the course id from question >subject > system, but that would
	//  cost resource, i rather add a courseID to the Question table and get it directly
	CourseID *uint   `json:"-"`
	Course   *Course `json:"course,omitempty" gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
}

//
type QuestionList struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}
type QuestionSearch struct {
	ID      uint   `json:"id"`
	Subject string `json:"subject"`
	System  string `json:"system"`
	Body    string `json:"body"`
	Course  string `json:"course"`
}
type AdminCreateQuestionInput struct {
	Question      string `json:"question" validate:"required"`
	Option1       string `json:"option1" validate:"required"`
	Option2       string `json:"option2" validate:"required"`
	Option3       string `json:"option3" validate:"required"`
	Option4       string `json:"option4" validate:"required"`
	CorrectOption uint   `json:"correct" validate:"required"`
	Description   string `json:"description" validate:"required"`
	SystemID      uint   `json:"systemID" validate:"required"`
}

// GORM HOOKS
func (u *Question) AfterFind(tx *gorm.DB) (err error) {
	if u.Image != nil {
		for i := 0; i < len(u.Image); i++ {
			// image exists
			imageURL := U.BaseURL + "/" + u.Image[i]
			u.Image[i] = imageURL
		}
	}
	return nil
}
