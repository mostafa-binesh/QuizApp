package models

import "encoding/json"

type Option struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Title string `json:"option,omitempty"`
	// Index is a alphabetic order thing, eg. first option of question 1's index is A, second is B and so on
	Index string `json:"index,omitempty"`
	// isCorrect is a bool thing, zero or one, i donno why i didn't set it as bool instead of uint
	IsCorrect uint `json:"status"`
	// relationships
	QuestionID *uint     `json:"-"`
	Question   *Question `json:"question,omitempty" gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	DropdownID *uint     `json:"-"`
	Dropdown   *Dropdown `json:"dropdown,omitempty" gorm:"foreignKey:DropdownID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
}
type AdminOptionInput struct {
	Title     string `form:"title" validate:"required"`
	IsCorrect bool   `form:"isCorrect" validate:"required"`
}

// FrontOption model has a dedicated MarshalJSON method
type FrontOption struct {
	ID            uint
	Title         string
	Index         string
	IsCorrectUint *uint
	IsCorrectBool *bool
}

// convert to json
func (fo FrontOption) MarshalJSON() ([]byte, error) {
	// If IsCorrectUint is not nil, include IsCorrectBool in the JSON output
	if fo.IsCorrectUint != nil {
		return json.Marshal(struct {
			ID            uint   `json:"id"`
			Title         string `json:"option,omitempty"`
			Index         string `json:"index,omitempty"`
			IsCorrectUint *uint  `json:"status,omitempty"`
		}{
			ID:            fo.ID,
			Title:         fo.Title,
			Index:         fo.Index,
			IsCorrectUint: fo.IsCorrectUint,
		})
	}

	// If IsCorrectUint is nil, omit IsCorrectBool from the JSON output
	return json.Marshal(struct {
		ID            uint   `json:"id"`
		Title         string `json:"option,omitempty"`
		Index         string `json:"index,omitempty"`
		IsCorrectBool *bool  `json:"status,omitempty"`
	}{
		ID:            fo.ID,
		Title:         fo.Title,
		Index:         fo.Index,
		IsCorrectBool: fo.IsCorrectBool,
	})
}
func ConvertOptionToFrontOption(options *[]Option, questionType QuestionType) *[]FrontOption {
	frontOptions := make([]FrontOption, len(*options))
	var isCorrectUint *uint
	var isCorrectBool *bool
	for i, option := range *options {
		isCorrectUint = nil
		isCorrectBool = nil
		// if question type was nextGeneration
		// and nextGeneration type was table(single or multiple)select, we need to set isCorrectUint
		// otherwise, we need to set isCorrectBool value
		if questionType == NextGenerationTableSingleSelect ||
			questionType == NextGenerationTableMultipleSelect {
			isCorrectUint = &option.IsCorrect
		} else {
			newBool := option.IsCorrect != 0
			isCorrectBool = &newBool
		}
		frontOptions[i] = FrontOption{
			ID:            option.ID,
			Title:         option.Title,
			Index:         option.Index,
			IsCorrectUint: isCorrectUint,
			IsCorrectBool: isCorrectBool,
		}
	}
	return &frontOptions
}
