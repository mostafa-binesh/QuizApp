package models

import "fmt"

const (
	MultipleSelect uint = iota
	SingleSelect
	NextGeneration
)

type Question struct {
	ID          uint    `json:"no" gorm:"primary_key"`
	Title       string  `json:"question"`
	Status      string  `json:"-"`
	Description string  `json:"description"`
	Images      []Image `gorm:"polymorphic:Owner;"`
	// relationships
	Options  []Option `json:"options,omitempty"`
	SystemID uint     `json:"-"`
	System   *System  `json:"system,omitempty" gorm:"foreignKey:SystemID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	// although we could get the course id from question >subject > system, but that would
	//  cost resource, i rather add a courseID to the Question table and get it directly
	CourseID *uint   `json:"-"`
	Course   *Course `json:"course,omitempty" gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	Type     uint    `json:"type"`
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
	Question     string             `json:"question" validate:"required"`
	Options      []AdminOptionInput `json:"options" validate:"required"`
	Description  string             `json:"description" validate:"required"`
	SystemID     uint               `json:"systemID" validate:"required"`
	QuestionType uint               `json:"questionType" validate:"required"`
}

func (question *Question) ConvertTypeStringToTypeInt(value string) error {
	switch value {
	case "multipleSelect":
		question.Type = MultipleSelect
	case "singleSelect":
		question.Type = SingleSelect
	case "nextGeneration":
		question.Type = NextGeneration
	default:
		return fmt.Errorf("Question type should be 'multipleSelect' or 'singleSelect or 'nextGeneration'")
	}
	return nil
}
