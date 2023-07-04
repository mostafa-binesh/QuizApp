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

func AddQuestion(c *fiber.Ctx) error {
	payload := new(M.AdminCreateQuestionInput)
	// parse body
	if err := c.BodyParser(payload); err != nil {
		return U.ResErr(c, err.Error())
	}
	// validate the payload
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	// init options array
	options := []*M.Option{}
	options = append(options, &M.Option{Title: payload.Option1, Index: "A"})
	options = append(options, &M.Option{Title: payload.Option2, Index: "B"})
	options = append(options, &M.Option{Title: payload.Option3, Index: "C"})
	options = append(options, &M.Option{Title: payload.Option4, Index: "D"})
	// if first option is correct, client needs to send 1
	options[payload.CorrectOption-1].IsCorrect = true
	img, _ := c.FormFile("image")
	var imgName *string
	if img != nil {
		uuid := uuid.New().String()
		newFileName := uuid + "-" + img.Filename
		c.SaveFile(img, fmt.Sprintf(U.UploadLocation+"/%s", newFileName))
		imgName = &newFileName
	}
	newQuestion := M.Question{
		Title:       payload.Question,
		Options:     options,
		SystemID:    payload.SystemID,
		Description: payload.Description,
		Image:       imgName,
	}
	result := D.DB().Create(&newQuestion)
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	return c.JSON(fiber.Map{"msg": "Question created successfully", "id": newQuestion.ID})
}
func QuestionByID(c *fiber.Ctx) error {
	question := &M.Question{}
	if err := D.DB().Preload("Course").Preload("System.Subject").First(question, c.Params("id")).Error; err != nil {
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
