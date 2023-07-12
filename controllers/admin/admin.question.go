package admin

import (
	D "docker/database"
	"fmt"
	"io"
	"os"

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
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	// handling options
	// init options array
	questionOptions := []M.Option{}
	questionOptions = append(questionOptions, M.Option{Title: payload.Option1, Index: "A"})
	questionOptions = append(questionOptions, M.Option{Title: payload.Option2, Index: "B"})
	questionOptions = append(questionOptions, M.Option{Title: payload.Option3, Index: "C"})
	questionOptions = append(questionOptions, M.Option{Title: payload.Option4, Index: "D"})
	// if first option is correct, client needs to send 1
	questionOptions[payload.CorrectOption-1].IsCorrect = true
	// # get images from request
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
func AdvancedUploadImage(c *fiber.Ctx) error {
	// Create a new file to store the uploaded data
	dst, err := os.Create("uploaded_file.txt")
	if err != nil {
		return err
	}
	defer dst.Close()

	// Get the request body stream
	reader := c.Request().BodyStream()

	// Read 1MiB at a time
	buffer := make([]byte, 0, 1024*1024)
	for {
		length, err := io.ReadFull(reader, buffer[:cap(buffer)])
		// Cap the buffer based on the actual length read
		buffer = buffer[:length]
		if err != nil {
			// When the error is EOF, there are no longer any bytes to read
			// meaning the request is completed
			if err == io.EOF {
				break
			}

			// If the error is an unexpected EOF, the requested size to read
			// was larger than what was available. This is not an issue for
			// as long as the length returned above is used, or the buffer
			// is capped as it is above. Any error that is not an unexpected
			// EOF is an actual error, which we handle accordingly
			if err != io.ErrUnexpectedEOF {
				return err
			}
		}

		// Write the buffered data to the destination file
		_, err = dst.Write(buffer)
		if err != nil {
			return err
		}
	}
	return c.SendString("DONE")
	// return c.JSON(fiber.Map{"data": fiber.Map{"img": c.BaseURL() + "/public/uploads/" + newFileName}})
}
