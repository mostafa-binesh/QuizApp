package controllers

import (
	D "docker/database"
	M "docker/models"
	U "docker/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AllQuestions(c *fiber.Ctx) error {
	// select all userAnswers of the authenticated user
	user := c.Locals("user").(M.User)
	if err := D.DB().Model(&user).
		Preload("Quizzes.UserAnswers.Question", func(db *gorm.DB) *gorm.DB { // could do this as well : Preload("Comments", "ORDER BY ? ASC > ?", "id")
			// if filterType is body and search query exists, search in question's title
			if c.Query("filterBy") == "body" && c.Query("search") != "" {
				db = db.Where("title ILIKE ?", fmt.Sprintf("%%%s%%", c.Query("search")))
			}
			return db
		}).
		Preload("Quizzes.UserAnswers.Question.System", func(db *gorm.DB) *gorm.DB { // could do this as well : Preload("Comments", "ORDER BY ? ASC > ?", "id")
			// if filterType is system and search query exists, search in subject's title
			if c.Query("filterBy") == "system" && c.Query("search") != "" {
				db = db.Where("title ILIKE ?", fmt.Sprintf("%%%s%%", c.Query("search")))
			}
			return db
		}).Preload("Quizzes.UserAnswers.Question.System.Subject", func(db *gorm.DB) *gorm.DB { // could do this as well : Preload("Comments", "ORDER BY ? ASC > ?", "id")
		// if filterType is subject and search query exists, search in question's title
		if c.Query("filterBy") == "subject" && c.Query("search") != "" {
			db = db.Where("title ILIKE ?", fmt.Sprintf("%%%s%%", c.Query("search")))
		}
		return db
	}).Preload("Quizzes.UserAnswers.Question.System.Subject.Course", func(db *gorm.DB) *gorm.DB { // could do this as well : Preload("Comments", "ORDER BY ? ASC > ?", "id")
		// if filterType is subject and search query exists, search in question's title
		if c.Query("filterBy") == "course" && c.Query("search") != "" {
			db = db.Where("title ILIKE ?", fmt.Sprintf("%%%s%%", c.Query("search")))
		}
		return db
	}).Find(&user).Error; err != nil {
		return U.DBError(c, err)
	}
	var questions []M.QuestionSearch
	for _, quiz := range user.Quizzes {
		for _, answer := range quiz.UserAnswers {
			if answer.Question != nil &&
				answer.Question.System != nil &&
				answer.Question.System.Subject != nil &&
				answer.Question.System.Subject.Course != nil {
				questions = append(questions, M.QuestionSearch{
					ID:      answer.Question.ID,
					System:  answer.Question.System.Title,
					Subject: answer.Question.System.Subject.Title,
					Course:  answer.Question.System.Subject.Course.Title,
					Body:    answer.Question.Title,
				})
			}
		}
	}
	return c.JSON(fiber.Map{"data": questions})
}
