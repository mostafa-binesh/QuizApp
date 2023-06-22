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
	user := c.Locals("user").(M.User)
	if err := D.DB().Model(&user).Preload("Quizzes.UserAnswers", func(db *gorm.DB) *gorm.DB { // could do this as well : Preload("Comments", "ORDER BY ? ASC > ?", "id")
		// if "title" query exists, search it in the database
		if c.Query("title") != "" {
			db = db.Where("Note LIKE ?", fmt.Sprintf("%%%s%%", c.Query("title")))
		}
		// select rows that Note field is not null
		db = db.Where("Note IS NOT NULL")
		return db
	}).Find(&user).Error; err != nil {
		return U.DBError(c, err)
	}
	// get notes only
	notes := []M.AnswerNote
	for _, quiz := range user.Quizzes {
		for _, answer := range quiz.UserAnswers {
			notes = append(notes, *answer.Note)
			notes = append(notes, M.AnswerNote{
				ID:         answer.ID,
				QuestionID: answer.QuestionID,
				Note:       answer.Note
			}
		}
	}
	return c.JSON(fiber.Map{"data": notes})
}
