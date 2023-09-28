package services

import (
	D "docker/database"
	M "docker/models"
	U "docker/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func SingleSelect(payload M.AdminCreateSingleSelectQuestionInput, c *fiber.Ctx) (*M.Question, error) {
	// handling options
	// init options arra
	questionOptions := []M.Option{}
	questionOptions = append(questionOptions, M.Option{Title: payload.Option1, Index: "A"})
	questionOptions = append(questionOptions, M.Option{Title: payload.Option2, Index: "B"})
	questionOptions = append(questionOptions, M.Option{Title: payload.Option3, Index: "C"})
	questionOptions = append(questionOptions, M.Option{Title: payload.Option4, Index: "D"})
	// if first option is correct, client needs to send 1
	questionOptions[payload.CorrectOption-1].IsCorrect = 1
	// # get images from request
	form, err := c.MultipartForm()
	if err != nil {
		return nil, U.ResErr(c, err.Error())
	}
	// images is optional
	images := form.File["images"]
	// create appropriate unique name for images and save them int disc
	var questionImages []M.Image
	for _, image := range images {
		uuid := uuid.New().String()
		newFileName := uuid + "-" + image.Filename
		c.SaveFile(image, fmt.Sprintf(U.UploadLocation+"/%s", newFileName))
		questionImages = append(questionImages, M.Image{
			Name: newFileName,
		})
	}
	var system M.System
	if err := D.DB().
		Preload("Subject").
		Find(&system).Error; err != nil {
		return nil, U.DBError(c, err)
	}
	// # create new question with given info
	newQuestion := M.Question{
		Title:       payload.Question,
		Options:     questionOptions,
		SystemID:    payload.SystemID,
		Description: payload.Description,
		Images:      questionImages,
		Type:        M.SingleSelect,
		CourseID:    &system.Subject.CourseID,
	}
	return &newQuestion, nil
}

func MultipleSelect(payload M.AdminCreateMultipleSelectQuestionInput, c *fiber.Ctx) (*M.Question, error) {
	// handling options
	questionOptions := []M.Option{}
	// go through each payload.option and convert it into M.Option
	for i, option := range payload.Options {
		questionOptions = append(questionOptions, M.Option{
			Title:     option.Title,
			Index:     U.GetNthAlphabeticUpperLetter(i + 1),
			IsCorrect: U.ConvertBoolToUint(option.IsCorrect),
		})
	}
	// get images from request
	form, err := c.MultipartForm()
	if err != nil {
		return nil, U.ResErr(c, "Parsing multipart form data failed")
	}
	// images := form.File["images"]
	// if images == nil {
	// 	return U.ResErr(c, err.Error())
	// }
	images := form.File["images"]
	// create appropriate unique name for images and save them int disc
	var questionImages []M.Image
	for _, image := range images {
		uuid := uuid.New().String()
		newFileName := uuid + "-" + image.Filename
		c.SaveFile(image, fmt.Sprintf(U.UploadLocation+"/%s", newFileName))
		questionImages = append(questionImages, M.Image{
			Name: newFileName,
		})
	}
	// create new question with given info
	var system M.System
	if err := D.DB().
		Preload("Subject").
		Find(&system).Error; err != nil {
		return nil, U.DBError(c, err)
	}
	newQuestion := M.Question{
		Title:       payload.Question,
		Options:     questionOptions,
		SystemID:    payload.SystemID,
		Description: payload.Description,
		Images:      questionImages,
		Type:        M.MultipleSelect,
		CourseID:    &system.Subject.CourseID,
	}
	// convert frontend's sent string question type to backend uint question type
	newQuestion.ConvertTypeStringToTypeInt(payload.QuestionType)
	return &newQuestion, nil
}
