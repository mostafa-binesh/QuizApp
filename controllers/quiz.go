package controllers

import (
	D "docker/database"
	"math/rand"
	"time"

	// F "docker/database/filters"
	M "docker/models"
	U "docker/utils"

	"github.com/gofiber/fiber/v2"
)

func AllQuizzes(c *fiber.Ctx) error {
	user := c.Locals("user").(M.User)
	result := D.DB().Model(&user).Preload("Quizzes").Find(&user)
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	userQuizzes := M.ConvertQuizToQuizList(user.Quizzes)
	return c.JSON(fiber.Map{"data": userQuizzes})
}

// ! CHECK: files ham preload mishe. aya niazi?
func QuizByID(c *fiber.Ctx) error {
	law := &M.Law{}
	if err := D.DB().Preload("Comments.User").Preload("Files").First(law, c.Params("id")).Error; err != nil {
		return U.DBError(c, err)
	}
	LawByID := M.LawToLawByID(law)
	return c.JSON(fiber.Map{
		"data": LawByID,
	})
}

func CreateQuiz(c *fiber.Ctx) error {
	payload := new(M.QuizInput)
	// parsing the payload
	if err := c.BodyParser(payload); err != nil {
		U.ResErr(c, err.Error())
	}
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	quiz := M.Quiz{
		UserID: 1,
		Status: "pending",
	}
	result := D.DB().Create(&quiz)
	if result.Error != nil {
		return U.ResErr(c, result.Error.Error())
	}
	var questionCount int64
	D.DB().Model(&M.Question{}).Count(&questionCount)
	for i := 0; i < payload.QuestionsCount; i++ {
		rand.Seed(time.Now().UnixNano())
		D.DB().Create(&M.UserAnswer{
			QuizID:     quiz.ID,
			QuestionID: uint(rand.Intn(int(questionCount))),
			IsMarked:   false,
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "مصوبه با موفقیت اضافه شد",
	})
}
