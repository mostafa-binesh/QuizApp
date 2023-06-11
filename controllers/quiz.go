package controllers

import (
	D "docker/database"
	"fmt"
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
	quiz := &M.Quiz{}
	// if err := D.DB().Preload("UserAnswers").Preload("UserAnswers.Question").First(quiz, c.Params("id")).Error; err != nil {
	if err := D.DB().Preload("UserAnswers.Question.Options").First(quiz, c.Params("id")).Error; err != nil {
		return U.DBError(c, err)
	}
	// LawByID := M.LawToLawByID(law)
	return c.JSON(fiber.Map{
		"data": quiz.ConvertQuizToQuizToFront(),
	})
}
func CreateQuiz(c *fiber.Ctx) error {
	payload := new(M.QuizInput)
	// parsing the payload
	if err := c.BodyParser(payload); err != nil {
		U.ResErr(c, err.Error())
	}
	fmt.Printf("payload: %v\n", payload)
	// validation the payload
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	user := c.Locals("user").(M.User)
	// create the quiz
	quiz := M.Quiz{
		UserID: user.ID,
		Status: "pending",
	}
	result := D.DB().Create(&quiz)
	if result.Error != nil {
		return U.ResErr(c, result.Error.Error())
	}
	// create the empty user answer h selected systems
	// TODO SECURITY ISSUE: checking the id of systems has not been checked
	var questionIDs []uint
	if err := D.DB().Model(&M.Question{}).Distinct().Where("system_id IN (?)", payload.SystemIDs).Pluck("id", &questionIDs).Error; err != nil {
		return U.DBError(c, err)
	}
	var questionCount int64
	D.DB().Model(&M.Question{}).Count(&questionCount)
	questionsCount := len(questionIDs)
	var randomIndex uint
	for i := 0; i < payload.QuestionsCount; i++ {
		randomIndex = uint(rand.Intn(int(questionsCount)))
		rand.Seed(time.Now().UnixNano())
		D.DB().Create(&M.UserAnswer{
			QuizID:     quiz.ID,
			QuestionID: questionIDs[randomIndex],
			IsMarked:   false,
			UserID:     user.ID,
			Status:     "unvisited",
		})
		U.RemoveElementByRef[uint](&questionIDs, int(randomIndex))
		questionsCount = len(questionIDs)
	}
	return c.Status(200).JSON(fiber.Map{
		"msg": "کوییز ساخته شد",
	})
}
