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
	"gorm.io/gorm"
)

func AllQuizzes(c *fiber.Ctx) error {
	user := c.Locals("user").(M.User)
	if err := D.DB().Model(&user).Preload("Quizzes.Course").Preload("UserAnswers", func(db *gorm.DB) *gorm.DB { // could do this as well : Preload("Comments", "ORDER BY ? ASC > ?", "id")
		db = db.Order("id ASC")
		return db
	}).Preload("Quizzes.UserAnswers.Question.Options").First(&user).Error; err != nil {
		return U.DBError(c, err)
	}
	var userQuizzes []M.QuizToFront
	for _, quiz := range user.Quizzes {
		fmt.Printf("course is %v\n", quiz.Course)
		userQuizzes = append(userQuizzes, quiz.ConvertQuizToQuizToFront())
	}
	return c.JSON(fiber.Map{"data": userQuizzes})
}

func QuizByID(c *fiber.Ctx) error {
	// get authenticated user
	user := c.Locals("user").(M.User)
	// find quiz with id of paramID of the user with dependencies
	if result := D.DB().Model(&user).Preload("Quizzes", c.Params("id")).Preload("Quizzes.Course").Preload("UserAnswers", func(db *gorm.DB) *gorm.DB { // could do this as well : Preload("Comments", "ORDER BY ? ASC > ?", "id")
		db = db.Order("id ASC")
		return db
	}).Preload("Quizzes.UserAnswers.Question.Options").First(&user); result.Error != nil {
		return U.DBError(c, result.Error)
	}
	// if user quiz with desired id doesn't exist
	fmt.Printf("user quizzes count: %d\n", len(user.Quizzes))
	if len(user.Quizzes) <= 0 {
		return U.ResErr(c, "Quiz not found")
	}
	return c.JSON(fiber.Map{"data": user.Quizzes[0].ConvertQuizToQuizToFront()})
}
func CreateQuiz(c *fiber.Ctx) error {
	payload := new(M.QuizInput)
	// parsing the payload
	if err := c.BodyParser(payload); err != nil {
		U.ResErr(c, err.Error())
	}
	// validation the payload
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	user := c.Locals("user").(M.User)
	// create the quiz
	endTime := time.Now().Add(time.Hour * 1)
	quiz := M.Quiz{
		UserID:   user.ID,
		Status:   "pending",
		EndTime:  &endTime,
		CourseID: 1, // TODO hardcoded !
	}
	result := D.DB().Create(&quiz)
	if result.Error != nil {
		return U.ResErr(c, result.Error.Error())
	}
	// create the empty user answer h selected systems
	// TODO SECURITY ISSUE: checking the id of systems has not been checked
	// all avaiable questions's id with desired system_id
	var questionIDs []uint
	if err := D.DB().Model(&M.Question{}).Distinct().Where("system_id IN (?)", payload.SystemIDs).Pluck("id", &questionIDs).Error; err != nil {
		return U.DBError(c, err)
	}
	questionsCount := len(questionIDs)
	var randomIndex uint
	fmt.Printf("payload question: %d , available questions count: %d\n", payload.QuestionsCount, questionsCount)
	if questionsCount <= payload.QuestionsCount {
		return U.ResErr(c, "not enough questions available")
	}
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
	fmt.Printf("printing the result\n")
	return c.Status(200).JSON(fiber.Map{
		"msg":    "Quiz been created",
		"quizID": quiz.ID,
	})
}
func UpdateQuiz(c *fiber.Ctx) error {
	payload := new(M.QuizToFront)
	// parsing the payload
	if err := c.BodyParser(payload); err != nil {
		U.ResErr(c, err.Error())
	}
	// validation the payload
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	// get authenticated user info
	user := c.Locals("user").(M.User)
	// ! todo maybe can optimize it
	// get quiz and its dependencies with paramID of the user in order
	// didn't get the error of ParamsInt, because i checked it in the router
	quizID, _ := c.ParamsInt("id")
	if err := D.DB().Model(&user).Preload("Quizzes", quizID).Preload("UserAnswers", func(db *gorm.DB) *gorm.DB { // could do this as well : Preload("Comments", "ORDER BY ? ASC > ?", "id")
		db = db.Order("id ASC")
		return db
	}).Preload("Quizzes.UserAnswers.Question.Options").First(&user).Error; err != nil {
		return U.DBError(c, err)
	}
	if !(len(user.Quizzes) == 1) {
		return U.ResErr(c, "This quiz doesn't exist")
	}
	quiz := M.Quiz{}
	quiz = user.Quizzes[0]
	convertedUserAnswer := payload.ConvertQuizFrontToQuiz(quiz.UserAnswers)
	// update the user answers into database
	fmt.Printf("convertedUserAnswer = %v\n", convertedUserAnswer)
	if err := D.DB().Save(convertedUserAnswer).Error; err != nil {
		return U.DBError(c, err)
	}
	return c.JSON(fiber.Map{"asd": convertedUserAnswer})
	return U.ResMessage(c, "Quiz been updated")
}
func CreateFakeQuiz(c *fiber.Ctx) error {
	payload := new(M.QuizInput)
	payload.QuestionsCount = 3
	payload.SystemIDs = []uint{1, 2, 3, 4, 5, 6, 7, 8, 9}
	// parsing the payload
	if err := c.BodyParser(payload); err != nil {
		U.ResErr(c, err.Error())
	}
	// validation the payload
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	user := c.Locals("user").(M.User)
	// create the quiz
	endTime := time.Now().Add(time.Hour * 1)
	quiz := M.Quiz{
		UserID:   user.ID,
		Status:   "pending",
		EndTime:  &endTime,
		CourseID: 1, // TODO hardcoded !
	}
	result := D.DB().Create(&quiz)
	if result.Error != nil {
		return U.ResErr(c, result.Error.Error())
	}
	// create the empty user answer h selected systems
	// TODO SECURITY ISSUE: checking the id of systems has not been checked
	// all avaiable questions's id with desired system_id
	var questionIDs []uint
	if err := D.DB().Model(&M.Question{}).Distinct().Where("system_id IN (?)", payload.SystemIDs).Pluck("id", &questionIDs).Error; err != nil {
		return U.DBError(c, err)
	}
	questionsCount := len(questionIDs)
	var randomIndex uint
	fmt.Printf("payload question: %d , available questions count: %d\n", payload.QuestionsCount, questionsCount)
	if questionsCount <= payload.QuestionsCount {
		return U.ResErr(c, "not enough questions available")
	}
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
	fmt.Printf("printing the result\n")
	return c.Status(200).JSON(fiber.Map{
		"msg":    "Quiz been created",
		"quizID": quiz.ID,
	})
}
