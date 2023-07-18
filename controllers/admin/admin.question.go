package admin

import (
	D "docker/database"
	"fmt"
	// "io"
	// "os"

	// F "docker/database/filters"
	M "docker/models"
	U "docker/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateMultipleSelectQuestion(c *fiber.Ctx) error {
	payload := new(M.AdminCreateMultipleSelectQuestionInput)
	// parse body
	if err := c.BodyParser(payload); err != nil {
		return U.ResErr(c, err.Error())
	}
	fmt.Printf("payload: %v\n", payload)
	// validate the payload
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
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
	fmt.Printf("after setting the options\n")
	// get images from request
	form, err := c.MultipartForm()
	if err != nil {
		return U.ResErr(c, "Parsing multipart form data failed")
	}
	fmt.Printf("after setting the multi part form\n")
	// images := form.File["images"]
	// if images == nil {
	// 	return U.ResErr(c, err.Error())
	// }
	images := form.File["images"]
	fmt.Printf("images variable %v\n", images)
	// create appropriate unique name for images and save them int disc
	var questionImages []M.Image
	for _, image := range images {
		fmt.Printf("image name: %v\n", image.Filename)
		uuid := uuid.New().String()
		newFileName := uuid + "-" + image.Filename
		c.SaveFile(image, fmt.Sprintf(U.UploadLocation+"/%s", newFileName))
		questionImages = append(questionImages, M.Image{
			Name: newFileName,
		})
	}
	// create new question with given info
	newQuestion := M.Question{
		Title:       payload.Question,
		Options:     questionOptions,
		SystemID:    payload.SystemID,
		Description: payload.Description,
		Images:      questionImages,
		Type:        M.MultipleSelect,
	}
	// convert frontend's sent string question type to backend uint question type
	newQuestion.ConvertTypeStringToTypeInt(payload.QuestionType)
	// insert new question to the database
	result := D.DB().Create(&newQuestion)
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	return c.JSON(fiber.Map{"msg": "Question created successfully", "id": newQuestion.ID})
}
func CreateSingleSelectQuestion(c *fiber.Ctx) error {
	payload := new(M.AdminCreateSingleSelectQuestionInput)
	// parse body
	if err := c.BodyParser(payload); err != nil {
		return U.ResErr(c, err.Error())
	}
	// validate the payload
	fmt.Printf("here 0")
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	// handling options
	// init options array
	fmt.Printf("here 1")
	questionOptions := []M.Option{}
	questionOptions = append(questionOptions, M.Option{Title: payload.Option1, Index: "A"})
	questionOptions = append(questionOptions, M.Option{Title: payload.Option2, Index: "B"})
	questionOptions = append(questionOptions, M.Option{Title: payload.Option3, Index: "C"})
	questionOptions = append(questionOptions, M.Option{Title: payload.Option4, Index: "D"})
	// if first option is correct, client needs to send 1
	questionOptions[payload.CorrectOption-1].IsCorrect = 1
	fmt.Printf("here 2")
	// # get images from request
	form, err := c.MultipartForm()
	if err != nil {
		return U.ResErr(c, err.Error())
	}
	fmt.Printf("here 2.5")
	images := form.File["images"]
	if images == nil {
		return U.ResErr(c, err.Error())
	}
	fmt.Printf("here 3")
	fmt.Printf("images variable %v\n", images)
	// create appropriate unique name for images and save them int disc
	var questionImages []M.Image
	for _, image := range images {
		fmt.Printf("image name: %v\n", image.Filename)
		uuid := uuid.New().String()
		newFileName := uuid + "-" + image.Filename
		c.SaveFile(image, fmt.Sprintf(U.UploadLocation+"/%s", newFileName))
		questionImages = append(questionImages, M.Image{
			Name: newFileName,
		})
	}
	fmt.Printf("here 4")
	// # create new question with given info
	newQuestion := M.Question{
		Title:       payload.Question,
		Options:     questionOptions,
		SystemID:    payload.SystemID,
		Description: payload.Description,
		Images:      questionImages,
		Type:        M.SingleSelect,
	}
	// # convert frontend's sent string question type to backend uint question type
	// newQuestion.ConvertTypeStringToTypeInt(payload.QuestionType)
	// insert new question to the database
	result := D.DB().Create(&newQuestion)
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	return c.JSON(fiber.Map{"msg": "Question created successfully", "id": newQuestion.ID})
}

// WIP
func CreateNextGenerationQuestion(c *fiber.Ctx) error {
	payload := new(M.AdminCreateNextGenerationQuestionInput)
	// parse body
	if err := c.BodyParser(payload); err != nil {
		return U.ResErr(c, err.Error())
	}
	// validate the payload
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	return c.JSON(fiber.Map{"data": payload})
}
func QuestionByID(c *fiber.Ctx) error {
	question := &M.Question{}
	// find the question with id of param id and preload course, iamges, system.subject
	if err := D.DB().Preload("Course").Preload("Images").Preload("System.Subject").First(question, c.Params("id")).Error; err != nil {
		return U.DBError(c, err)
	}
	return c.JSON(fiber.Map{"data": question})
}

func UploadImage(c *fiber.Ctx) error {
	type Upload struct {
		File string `validate:"required"`
	}
	payload := new(Upload)
	file, err := c.FormFile("file")
	if file != nil {
		payload.File = file.Filename
	}
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	// check if file with this name already exists
	if U.FileExistenceCheck(file.Filename, U.UploadLocation) {
		return U.ResErr(c, "file already exists")
	}
	// ! file extension check
	fmt.Printf("file name %s\n", file.Filename)
	fmt.Printf("extension check %t\n", U.IsImageFile(file.Filename))
	if !(U.IsImageFile(file.Filename)) {
		return c.SendString("file should be image! please fix it")
	}
	// Save file to disk
	uuid := uuid.New().String()
	newFileName := uuid + "-" + file.Filename
	err = c.SaveFile(file, fmt.Sprintf(U.UploadLocation+"/%s", newFileName))
	if err != nil {
		return U.ResErr(c, "cannot save | "+err.Error())
	}
	return c.JSON(fiber.Map{"data": fiber.Map{"img": c.BaseURL() + "/public/uploads/" + newFileName}})
}
func ChangeImageURLsInDescription(c *fiber.Ctx) error {
	type ChangeImageURL struct {
		PreviousWebsite string `json:"previousWebsite"`
		NewWebsite      string `json:"newWebsite"`
	}
	payload := new(ChangeImageURL)
	// parse body
	if err := c.BodyParser(payload); err != nil {
		return U.ResErr(c, err.Error())
	}
	// get all questions
	var questions []M.Question
	if err := D.DB().Find(&questions).Error; err != nil {
		return U.DBError(c, err)
	}
	// replace the previous website url with the new one
	for i, _ := range questions {
		questions[i].ChangeImageURLsInDescription(payload.PreviousWebsite, payload.NewWebsite)
	}
	// save modified questions
	if err := D.DB().Save(&questions).Error; err != nil {
		return U.DBError(c, err)
	}
	return U.ResMessage(c, "Success")
}
