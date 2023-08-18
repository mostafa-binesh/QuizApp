package controllers

import (
	D "docker/database"
	M "docker/models"
	U "docker/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AllNotes(c *fiber.Ctx) error {
	// select all userAnswers of the authenticated user
	user := M.AuthedUser(c)
	if err := D.DB().Model(&user).Preload("Quizzes", func(db *gorm.DB) *gorm.DB { // could do this as well : Preload("Comments", "ORDER BY ? ASC > ?", "id")
		// if "quizID" query exists, search it in the database
		if c.Query("quizID") != "" {
			db = db.Where(c.QueryInt("quizID"))
		}
		return db
	}).Preload("Quizzes.UserAnswers",
		func(db *gorm.DB) *gorm.DB { // could do this as well : Preload("Comments", "ORDER BY ? ASC > ?", "id")
			// if "title" query exists, search it in the database
			if c.Query("body") != "" {
				db = db.Where("Note ILIKE ?", fmt.Sprintf("%%%s%%", c.Query("body")))
			}
			if c.Query("sort") != "" {
				switch c.Query("sort") {
				case "newer":
					db = db.Order("id desc")
				case "older":
					db = db.Order("id asc")
				case "questionID":
					db = db.Order("QuestionID asc")
				}
			}
			// required: select rows that Note field is not null
			db = db.Where("Note IS NOT NULL")
			return db
		}).Preload("Quizzes.UserAnswers.Question").Find(&user).Error; err != nil {
		return U.DBError(c, err)
	}
	// get notes only
	var notes []M.AnswerNote
	for _, quiz := range user.Quizzes {
		for _, answer := range quiz.UserAnswers {
			notes = append(notes, M.AnswerNote{
				ID:       answer.ID,
				Question: answer.Question,
				Note:     answer.Note,
				QuizID:   answer.QuizID,
			})
		}
	}
	return c.JSON(fiber.Map{"data": notes})
}
func EditNote(c *fiber.Ctx) error {
	payload := new(M.EditNoteInput)
	// parsing the payload
	if err := c.BodyParser(payload); err != nil {
		U.ResErr(c, err.Error())
	}
	// validation the payload
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	// get userAnswer with url id param
	userAnswer := M.UserAnswer{}
	if err := D.DB().First(&userAnswer, c.Params("id")).Error; err != nil {
		return U.DBError(c, err)
	}
	userAnswer.Note = payload.Note
	if err := D.DB().Save(&userAnswer).Error; err != nil {
		return U.DBError(c, err)
	}
	return U.ResMessage(c, "Note has been updated")
}
