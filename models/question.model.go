package models

import (
	U "docker/utils"
	"fmt"
	"strings"
)

type QuestionType int8

const (
	MultipleSelect QuestionType = iota
	SingleSelect
	NextGeneration
	NextGenerationMultipleSelect
	NextGenerationSingleSelect
	NextGenerationTableSingleSelect
	NextGenerationTableMultipleSelect
	NextGenerationTableDropDown
)

type Question struct {
	ID          uint         `json:"no" gorm:"primary_key"`
	Title       string       `json:"question"`
	Status      string       `json:"-"`
	Description string       `json:"description"`
	Type        QuestionType `json:"type"`
	// # relationships
	CourseID    *uint        `json:"-"` // although we could get the course id from question >subject > system, but that would,  cost resource, i rather add a courseID to the Question table and get it directly
	Course      *Course      `json:"course,omitempty" gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	SystemID    uint         `json:"-"`
	System      *System      `json:"system,omitempty" gorm:"foreignKey:SystemID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	Images      []Image      `gorm:"polymorphic:Owner;"`
	Options     []Option     `json:"options,omitempty" gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"` // this cascade means if this question is deleted, all options will be deleted too and you won't encounter dependency error
	Tabs        []Tab        `json:"tabs"`
	Dropdowns   []Dropdown   `json:"dropdowns,omitempty"`
	UserAnswers []UserAnswer `json:"-"`
}
type FrontQuestion struct {
	ID          uint    `json:"no" gorm:"primary_key"`
	Title       string  `json:"question"`
	Description string  `json:"description"`
	Images      []Image `gorm:"polymorphic:Owner;"`
	// relationships
	Options  []FrontOption `json:"options,omitempty"`
	SystemID uint          `json:"-"`
	System   *System       `json:"system,omitempty" gorm:"foreignKey:SystemID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	// although we could get the course id from question >subject > system, but that would
	//  cost resource, i rather add a courseID to the Question table and get it directly
	Course                   *Course      `json:"course,omitempty" gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	Type                     QuestionType `json:"type"`
	Tabs                     []Tab        `json:"tabs"`
	Dropdowns                []Dropdown   `json:"dropdowns,omitempty"`
	AnswerAccuracyPercentage uint         `json:"answerAccuracyPercentage"`
}

func (question Question) ToFrontQuestion() FrontQuestion {
	frontQuestion := FrontQuestion{
		ID:                       question.ID,
		Title:                    question.Title,
		Description:              question.Description,
		Images:                   question.Images,
		Options:                  *ConvertOptionToFrontOption(&question.Options, question.Type),
		SystemID:                 question.SystemID,
		System:                   question.System,
		Course:                   question.Course,
		Type:                     question.Type,
		Tabs:                     question.Tabs,
		Dropdowns:                question.Dropdowns,
		AnswerAccuracyPercentage: uint(question.AnswerAccuracyPercentage()),
	}
	return frontQuestion
}
func ConvertQuestionsToFrontQuestions(questions *[]Question) *[]FrontQuestion {
	frontQuestions := make([]FrontQuestion, len(*questions))
	for i, question := range *questions {
		frontQuestions[i] = question.ToFrontQuestion()
	}
	return &frontQuestions
}

// question.UserAnswers must be preloaded
func (question Question) AnswerAccuracyPercentage() int {
	var correctAnswersCount int
	totalAnswersCount := len(question.UserAnswers)
	// Count the number of correct answers directly in the database query
	for _, answer := range question.UserAnswers {
		if answer.IsCorrect != nil && *answer.IsCorrect {
			correctAnswersCount++
		}
	}
	if totalAnswersCount == 0 {
		// Handle the case where there are no user answers to avoid division by zero
		return 0
	}
	// Calculate the accuracy percentage
	accuracyPercentage := (correctAnswersCount * 100) / totalAnswersCount
	return accuracyPercentage
}

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
type AdminCreateMultipleSelectQuestionInput struct {
	Question     string             `form:"question" validate:"required"`
	Options      []AdminOptionInput `form:"options" validate:"required"`
	Description  string             `form:"description" validate:"required"`
	SystemID     uint               `form:"systemID" validate:"required"`
	QuestionType string             `form:"questionType" validate:"required"`
}
type AdminCreateSingleSelectQuestionInput struct {
	Question      string `json:"email" validate:"required"`
	Option1       string `json:"option1" validate:"required"`
	Option2       string `json:"option2" validate:"required"`
	Option3       string `json:"option3" validate:"required"`
	Option4       string `json:"option4" validate:"required"`
	CorrectOption uint   `json:"correct" validate:"required"`
	Description   string `json:"description" validate:"required"`
	SystemID      uint   `json:"systemID" validate:"required"`
}
type AdminCreateNextGenerationQuestionInput struct {
	Title       string                                       `json:"question" validate:"required"`
	Description string                                       `json:"metaDescription" validate:"required"`
	Tabs        []AdminCreateNextGenerationQuestionTabsInput `json:"tabs" validate:"required"`
	Type        string                                       `json:"type" validate:"required"`
	// the options is for dropdown
	DropDownOptions [][]AdminCreateNextGenerationQuestionOptionsInput `json:"dropDownOptions"`
	SingleOptions   AdminCreateNextGenerationQuestionOptionsInput     `json:"singleOptions"`
	MultipleOptions []AdminCreateNextGenerationQuestionOptionsInput   `json:"multipleOptions"`
}

func (payload AdminCreateNextGenerationQuestionInput) ToQuestion() (q Question) {
	q.Title = payload.Title
	q.Description = payload.Description
	// tabs
	if err := U.Convert[AdminCreateNextGenerationQuestionTabsInput, Tab](payload.Tabs, &q.Tabs); err != nil {
		panic("conversion error in toQuestion of AdminCreateNextGenerationQuestionInput")
	}
	// type
	q.ConvertTypeStringToTypeInt(payload.Type)
	// dropdown options
	// ! needs attention
	if err := U.Convert[AdminCreateNextGenerationQuestionOptionsInput, Dropdown](payload.DropDownOptions[1], &q.Dropdowns); err != nil {
		panic("conversion error in toQuestion of AdminCreateNextGenerationQuestionInput")
	}
	return
}

// each tab has a name and tables
type AdminCreateNextGenerationQuestionTabsInput struct {
	TableTitle string     `json:"tableTitle" validate:"required"`
	Rows       [][]string `json:"rows" validate:"required"`
}
type AdminCreateNextGenerationQuestionOptionsInput struct {
	Title  string `json:"option"`
	Status any    `json:"status"`
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
func (question *Question) ConvertNextGenerationTypeToTypeInt(value string) error {
	switch value {
	case "singleSelect":
		question.Type = NextGenerationSingleSelect
	case "multipleSelect":
		question.Type = NextGenerationMultipleSelect
	case "tableSingleSelect":
		question.Type = NextGenerationTableSingleSelect
	case "tableMultipleSelect":
		question.Type = NextGenerationTableMultipleSelect
	case "dropdown":
		question.Type = NextGenerationTableDropDown
	default:
		return fmt.Errorf("Question type should be singleSelect or 'multipleSelect or TableSingleSelect or TableMultipleSelect or TableDropDown")
	}
	return nil
}

// replace previous webiste url with new website in description field
func (question *Question) ReplacePreWebsiteWithNewWebsiteImageURLDescription(previousSite string, newSite string) {
	question.Description = strings.Replace(question.Description, previousSite, newSite, -1)
}

// question.options must be preloaded
func (question Question) CorrectOptionsCount() (correctAnswersCount uint) {
	for _, option := range question.Options {
		if option.IsCorrect == 1 {
			correctAnswersCount++
		}
	}
	return
}
func (q Question) IsTraditional() bool {
	return q.Type == MultipleSelect || q.Type == SingleSelect
}

// return if the question is single select, despite the traditional or next generation
func (q Question) IsSingleSelect() bool {
	return q.Type == SingleSelect || q.Type == NextGenerationSingleSelect || q.Type == NextGenerationTableSingleSelect
}
