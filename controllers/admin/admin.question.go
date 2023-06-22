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
	options[payload.CorrectOption-1].IsCorrect = true
	newQuestion := M.Question{
		Title:       payload.Question,
		Options:     options,
		SystemID:    payload.SystemID,
		Description: payload.Description,
	}
	result := D.DB().Create(&newQuestion)
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	return U.ResMessage(c, "Question has been created")
}
func UploadImage(c *fiber.Ctx) error {
	type Upload struct {
		File string `json:"file" validate:"required"`
	}
	payload := new(Upload)
	if err := c.BodyParser(payload); err != nil {
		return c.JSON(fiber.Map{
			"error": err,
		})
	}
	file, err := c.FormFile("file")
	// if err != nil {
	// ! if file not exists, we get error: there is no uploaded file associated with the given key
	// 	return c.JSON(fiber.Map{"error": err.Error()})
	// }
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
