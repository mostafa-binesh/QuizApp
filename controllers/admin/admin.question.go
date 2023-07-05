package admin

import (
	D "docker/database"
	"fmt"

	// F "docker/database/filters"
	M "docker/models"
	U "docker/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateQuestion(c *fiber.Ctx) error {
	payload := new(M.AdminCreateQuestionInput)
	// parse body
	if err := c.BodyParser(payload); err != nil {
		return U.ResErr(c, err.Error())
	}
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
			IsCorrect: option.IsCorrect,
		}) 
	}
	// get images from request
	form, err := c.MultipartForm()
	images := form.File["images"]
	if err != nil || images == nil {
		return U.ResErr(c, err.Error())
	}
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
	}
	// insert new question to the database
	result := D.DB().Create(&newQuestion)
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	return c.JSON(fiber.Map{"msg": "Question created successfully", "id": newQuestion.ID})
}
func QuestionByID(c *fiber.Ctx) error {
	question := &M.Question{}
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
