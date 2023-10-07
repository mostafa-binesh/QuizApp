package admin

import (
	D "docker/database"
	"fmt"

	// "io"
	// "os"

	// F "docker/database/filters"
	F "docker/database/filters"
	M "docker/models"
	S "docker/services"
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
	newQuestion, err := S.MultipleSelect(*payload, c)
	if err != nil {
		return err
	}
	// insert new question to the database
	result := D.DB().Create(newQuestion)
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	return c.JSON(fiber.Map{"msg": "Question created successfully", "id": newQuestion.ID})
}
func EditMultipleSelectQuestion(c *fiber.Ctx) error {
	payload := new(M.AdminCreateMultipleSelectQuestionInput)
	// parse body
	if err := c.BodyParser(payload); err != nil {
		return U.ResErr(c, err.Error())
	}
	// validate the payload
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	// use the service to convert payload to question model
	editedQuestion, err := S.MultipleSelect(*payload, c)
	if err != nil {
		return err
	}
	// because it's about the update, set the id from query param
	// ignore the paramsInt error, because i've checked it already in the router
	questionID, _ := c.ParamsInt("questionID")
	editedQuestion.ID = uint(questionID)
	// delete existing options associated with the question
	if err := D.DB().Where("question_id = ?", questionID).Delete(M.Option{}).Error; err != nil {
		return err
	}
	// insert new question to the database
	result := D.DB().Save(editedQuestion)
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	return c.JSON(fiber.Map{"msg": "Question was updated"})
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

	// # create new question with given info
	newQuestion, err := S.SingleSelect(*payload, c)
	if err != nil {
		return err
	}
	// # convert frontend's sent string question type to backend uint question type
	// newQuestion.ConvertTypeStringToTypeInt(payload.QuestionType)
	// insert new question to the database
	result := D.DB().Create(newQuestion)
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	return c.JSON(fiber.Map{"msg": "Question created successfully", "id": newQuestion.ID})
}
func EditSingleSelectQuestion(c *fiber.Ctx) error {
	payload := new(M.AdminCreateSingleSelectQuestionInput)
	// parse body
	if err := c.BodyParser(payload); err != nil {
		return U.ResErr(c, err.Error())
	}
	// validate the payload
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}

	// # create new question with given info
	editedQuestion, err := S.SingleSelect(*payload, c)
	if err != nil {
		// return the error only, because we handlede the rest of it the in the service
		return err
	}
	// ignore the paramsInt error, because i've checked it already in the router
	questionID, _ := c.ParamsInt("questionID")
	editedQuestion.ID = uint(questionID)
	// delete existing options associated with the question
	if err := D.DB().Where("question_id = ?", questionID).Delete(M.Option{}).Error; err != nil {
		return err
	}
	// insert new question to the database
	result := D.DB().Save(editedQuestion)
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	return c.JSON(fiber.Map{"msg": "Question was updated"})
}

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
	if err := D.DB().
		Preload("Course").
		Preload("Images").
		Preload("System.Subject").
		Preload("Options").
		Preload("Tabs").
		Preload("Dropdowns").
		First(question, c.Params("id")).Error; err != nil {
		return U.DBError(c, err)
	}
	return c.JSON(fiber.Map{"data": question})
}

// TODO this handler and other upload handlers have memory leak !
func UploadImage(c *fiber.Ctx) error {
	type Upload struct {
		File string `validate:"required"`
		// FileName string `validate:"required"`
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
	if !U.IsImageFile(file.Filename) {
		return c.SendString("file should be image! please fix it")
	}
	// Save file to disk
	uuid := uuid.New().String()
	newFileName := uuid + "-" + file.Filename
	if err = c.SaveFile(file, fmt.Sprintf(U.UploadLocation+"/%s", newFileName)); err != nil {
		return U.ResErr(c, "cannot save | "+err.Error())
	}
	return c.JSON(fiber.Map{"data": fiber.Map{"img": c.BaseURL() + "/public/uploads/" + newFileName}})
}

// this handler and route has been created for when the creator changed the url of the website
// everything will be moved there, but image urls that been used in description of questions won't change
// this route will find and replace old website url with new website
func ChangeImageURLsInDescription(c *fiber.Ctx) error {
	type ChangeImageURL struct {
		PreviousWebsite string `json:"previousWebsite" validate:"required"`
		NewWebsite      string `json:"newWebsite" validate:"required"`
	}
	payload := new(ChangeImageURL)
	// parse body
	if err := c.BodyParser(payload); err != nil {
		return U.ResErr(c, err.Error())
	}
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	// get all questions
	var questions []M.Question
	if err := D.DB().Find(&questions).Error; err != nil {
		return U.DBError(c, err)
	}
	// replace the previous website url with the new one
	for i, _ := range questions {
		questions[i].ReplacePreWebsiteWithNewWebsiteImageURLDescription(payload.PreviousWebsite, payload.NewWebsite)
	}
	// save modified questions
	if err := D.DB().Save(&questions).Error; err != nil {
		return U.DBError(c, err)
	}
	return U.ResMsg(c, "Success")
}
func AllQuestions(c *fiber.Ctx) error {
	// get all questions
	var questions []M.Question
	// two available filter: type and courseID
	if err := D.DB().
		Scopes(
			F.FilterByType(c, F.FilterType{QueryName: "type"},
				F.FilterType{QueryName: "courseID", ColumnName: "course_id"})).
		Find(&questions).Error; err != nil {
		return U.DBError(c, err)
	}
	return c.JSON(fiber.Map{"data": questions})
}
