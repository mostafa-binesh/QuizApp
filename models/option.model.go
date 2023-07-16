package models

type Option struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	Title     string `json:"option,omitempty"`
	Index     string `json:"index,omitempty"`
	IsCorrect uint   `json:"status"`
	// relationships
	QuestionID uint      `json:"-"`
	Question   *Question `json:"question,omitempty" gorm:"foreignKey:QuestionID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
}
type AdminOptionInput struct {
	Title     string `form:"title" validate:"required"`
	IsCorrect bool   `form:"isCorrect" validate:"required"`
}
type FrontOption struct {
	ID            uint   `json:"id" gorm:"primary_key"`
	Title         string `json:"option,omitempty"`
	Index         string `json:"index,omitempty"`
	IsCorrectUint *uint  `json:"status,omitempty"`
	IsCorrectBool *bool  `json:"status,omitempty"`
}

func ConvertOptionToFrontOption(options *[]Option, questionType QuestionType) *[]FrontOption {
	frontOptions := make([]FrontOption, len(*options))
	var isCorrectUint *uint
	var isCorrectBool *bool
	for _, option := range *options {
		isCorrectUint = nil
		isCorrectBool = nil
		// if question type was nextGeneration
		// and nextGeneration type was table(single or multiple)select, we need to set isCorrectUint
		// otherwise, we need to set isCorrectBool value
		if questionType == NextGenerationTableSingleSelect ||
			questionType == NextGenerationTableMultipleSelect {
			*isCorrectUint = option.IsCorrect
		} else {
			*isCorrectBool = option.IsCorrect != 0
		}
		frontOptions = append(frontOptions, FrontOption{
			ID:            option.ID,
			Title:         option.Title,
			Index:         option.Index,
			IsCorrectUint: isCorrectUint,
			IsCorrectBool: isCorrectBool,
		})
	}
	return &frontOptions
}
