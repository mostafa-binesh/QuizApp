package controllers

import (
	D "docker/database"
	"fmt"
	"math/rand"
	"strings"
	"time"

	// F "docker/database/filters"
	M "docker/models"
	U "docker/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// because frontend guys only request once for all of quizzes and
// save the entire quizzes into the state, i need to send quiz with all of its info
func AllQuizzes(c *fiber.Ctx) error {
	user := c.Locals("user").(M.User)
	if err := D.DB().Model(&user).Preload("Quizzes.Course").Preload("UserAnswers", func(db *gorm.DB) *gorm.DB { // could do this as well : Preload("Comments", "ORDER BY ? ASC > ?", "id")
		db = db.Order("id ASC")
		return db
	}).Preload("Quizzes.UserAnswers.Question.Options").Preload("Quizzes.UserAnswers.Question.Dropdowns.Options").First(&user).Error; err != nil {
		return U.DBError(c, err)
	}
	var userQuizzes []M.QuizToFront
	for _, quiz := range user.Quizzes {
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
	}).Preload("Quizzes.UserAnswers.Question.Options").Preload("Quizzes.UserAnswers.Question.Dropdowns.Options").First(&user); result.Error != nil {
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
	// get course id using first system id
	if len(payload.SystemIDs) == 0 {
		return U.ResErr(c, "You must at least select one system")
	}
	// find the system from first index of systemIDs
	// cause we checked the length of payload.systemIDs > 0, we can use first index of it
	systemID := payload.SystemIDs[0]
	system := M.System{}
	D.DB().Preload("Subject.Course").Find(&system, systemID)
	// create the quiz
	endTime := time.Now().Add(time.Hour * 1) // todo: hardcoded, should get it from woocommerce.course.duration
	currentTime := time.Now()
	duration := endTime.Sub(currentTime)
	remainingSeconds := uint(duration.Seconds())
	// # create quiz
	quiz := M.Quiz{
		UserID:       user.ID,
		Status:       "pending",
		EndTime:      &endTime,
		Duration:     remainingSeconds,
		CourseID:     system.Subject.CourseID,
		Mode:         strings.Join(payload.QuizMode, ","),
		Type:         strings.Join(payload.QuizType, ","),
		QuestionMode: payload.QuestionMode,
	}
	result := D.DB().Create(&quiz)
	if result.Error != nil {
		return U.ResErr(c, result.Error.Error())
	}
	// # find questions based on question mode
	// create the empty user answer h selected systems
	// TODO SECURITY ISSUE: checking the id of systems has not been checked
	// all avaiable questions's id with desired system_id
	var questionIDs []uint
	var err error
	if payload.QuestionMode == M.AllQuestionMode {
		err = D.DB().Model(&M.Question{}).Distinct().Where("system_id IN (?)", payload.SystemIDs).Pluck("id", &questionIDs).Error
	} else if payload.QuestionMode == M.MarkedQuestionMode {
		err = D.DB().Model(&user).Preload("UserAnswers", "is_marked = ?", true).Find(&user).Error
		if err != nil {
			D.DB().Delete(&quiz)
			return U.DBError(c, err)
		}
		for _, answer := range user.UserAnswers {
			// we don't want duplicate questions, check if the questionID is already exist or not
			if !U.ExistsInArray[uint](questionIDs, answer.QuestionID) {
				questionIDs = append(questionIDs, answer.QuestionID)
			}
		}
	} else if payload.QuestionMode == M.SingleSelectQuestionMode {
		err = D.DB().Model(&M.Question{}).Where("type = ?", M.SingleSelect).Pluck("id", &questionIDs).Error
	} else if payload.QuestionMode == M.MultipleSelectQuestionMode {
		err = D.DB().Model(&M.Question{}).Where("type = ?", M.MultipleSelect).Pluck("id", &questionIDs).Error
	}
	if err != nil {
		D.DB().Delete(&quiz)
		return U.DBError(c, err)
	}
	questionsCount := len(questionIDs)
	var randomIndex uint
	fmt.Printf("payload question: %d , available questions count: %d\n", payload.QuestionsCount, questionsCount)
	if questionsCount <= payload.QuestionsCount {
		D.DB().Delete(&quiz)
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
	fmt.Printf("user's quizzes count: %d\n", len(user.Quizzes))
	fmt.Printf("user's quizzes: %v\n", user.Quizzes)
	quiz := M.Quiz{}
	for _, q := range user.Quizzes {
		// return c.JSON(fiber.Map{"data": q})
		quiz = q
		break
	}
	if len(user.Quizzes) == 0 {
		return U.ResErr(c, "This quiz doesn't exist")
	}
	if user.Quizzes[0].Status == "" {
		return U.ResErr(c, "This quiz doesn't exist")
	}
	// quiz = user.Quizzes[0]
	convertedUserAnswer := payload.ConvertQuizFrontToQuiz(quiz.UserAnswers)
	// update the user answers into database
	fmt.Printf("convertedUserAnswer = %v\n", convertedUserAnswer)
	if err := D.DB().Save(convertedUserAnswer).Error; err != nil {
		return U.DBError(c, err)
	}
	quiz.Status = payload.QuizState
	quiz.CalculateRemainingSeconds(payload.RemainingHours, payload.RemainingMinutes, payload.RemainingSeconds)
	if err := D.DB().Save(&quiz).Error; err != nil {
		return U.DBError(c, err)
	}
	// todo not tested yet
	// return c.JSON(fiber.Map{"asd": convertedUserAnswer})
	return U.ResMessage(c, "Quiz been updated")
}
func CreateFakeQuiz(c *fiber.Ctx) error {
	payload := new(M.QuizInput)
	payload.QuestionsCount = 3
	payload.SystemIDs = []uint{1, 2, 3, 4, 5, 6, 7, 8, 9}
	payload.QuizMode = []string{"tutor", "timed"}
	payload.QuizType = []string{"nextGeneration"}
	payload.QuestionMode = 1
	// parsing the payload
	if err := c.BodyParser(payload); err != nil {
		U.ResErr(c, err.Error())
	}
	// validation the payload
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	user := c.Locals("user").(M.User)

	if len(payload.SystemIDs) == 0 {
		return U.ResErr(c, "You must at least select one system")
	}
	// systemID := payload.SystemIDs[0]
	system := M.System{}
	if err := D.DB().Preload("Subject.Course").First(&system).Error; err != nil {
		return U.DBError(c, err)
	}
	// create the quiz
	endTime := time.Now().Add(time.Hour * 1)
	currentTime := time.Now()
	duration := endTime.Sub(currentTime)
	remainingSeconds := uint(duration.Seconds())
	quiz := M.Quiz{
		UserID:   user.ID,
		Status:   "pending",
		EndTime:  &endTime,
		Duration: remainingSeconds,
		CourseID: system.Subject.CourseID,
		Mode:     strings.Join(payload.QuizMode, ","),
		Type:     strings.Join(payload.QuizType, ","),
	}
	result := D.DB().Create(&quiz)
	if result.Error != nil {
		return U.ResErr(c, result.Error.Error())
	}
	// create the empty user answer h selected systems
	// TODO SECURITY ISSUE: checking the id of systems has not been checked
	// all avaiable questions's id with desired system_id
	var multipleSelectQuestionIDs []uint
	var singleSelectQuestionIDs []uint
	// if err := D.DB().Model(&M.Question{}).Distinct().Where("system_id IN (?)", payload.SystemIDs).Pluck("id", &questionIDs).Error; err != nil {
	// 	return U.DBError(c, err)
	// }
	// get 3 multiple question type
	if err := D.DB().Model(&M.Question{}).Where("type = ?", M.MultipleSelect).Limit(3).Pluck("id", &multipleSelectQuestionIDs).Error; err != nil {
		return U.DBError(c, err)
	}
	if err := D.DB().Model(&M.Question{}).Where("type = ?", M.SingleSelect).Limit(3).Pluck("id", &singleSelectQuestionIDs).Error; err != nil {
		return U.DBError(c, err)
	}
	multipleSelectQuestionIDs = append(multipleSelectQuestionIDs, singleSelectQuestionIDs...)
	questionsCount := len(multipleSelectQuestionIDs)
	var randomIndex uint
	fmt.Printf("payload question: %d , available questions count: %d\n", payload.QuestionsCount, questionsCount)
	if questionsCount < payload.QuestionsCount {
		D.DB().Delete(&quiz)
		return U.ResErr(c, "not enough questions available")
	}
	for i := 0; i < payload.QuestionsCount; i++ {
		randomIndex = uint(rand.Intn(int(questionsCount)))
		rand.Seed(time.Now().UnixNano())
		D.DB().Create(&M.UserAnswer{
			QuizID:     quiz.ID,
			QuestionID: multipleSelectQuestionIDs[randomIndex],
			IsMarked:   false,
			UserID:     user.ID,
			Status:     "unvisited",
		})
		U.RemoveElementByRef[uint](&multipleSelectQuestionIDs, int(randomIndex))
		questionsCount = len(multipleSelectQuestionIDs)
	}
	fmt.Printf("printing the result\n")
	return c.Status(200).JSON(fiber.Map{
		"msg":    "Quiz been created",
		"quizID": quiz.ID,
	})
}
func OverallReport(c *fiber.Ctx) error {
	user := c.Locals("user").(M.User)
	// options is needed in user preload in correct and incorrect answer coount
	if err := D.DB().Preload("Quizzes.UserAnswers.Question.Options").Preload("Courses.Subjects.Systems.Questions").Find(&user).Error; err != nil {
		return U.DBError(c, err)
	}
	var totalQuestionsCount int
	// todo: if we add course id to question model, here would be better solution
	for _, course := range user.Courses {
		for _, subject := range course.Subjects {
			for _, system := range subject.Systems {
				totalQuestionsCount += len(system.Questions)
			}
		}
	}
	var usedQuestions []uint
	// we need to get all unique questions
	for _, quiz := range user.Quizzes {
		for _, answer := range quiz.UserAnswers {
			if answer.Question != nil {
				// answer.Question.System != nil &&
				// answer.Question.System.Subject != nil &&
				// answer.Question.System.Subject.Course != nil {
				if !U.ExistsInArray[uint](usedQuestions, answer.Question.ID) {
					usedQuestions = append(usedQuestions, answer.Question.ID)
				}
			}
		}
	}
	usedQuestionsCount := len(usedQuestions)
	createdTests := len(user.Quizzes)
	var finishedTests int
	var suspendedTests int
	var userAnswers []M.UserAnswer
	for _, quiz := range user.Quizzes {
		if quiz.Status == "finished" { // todo: inha (finsihed, pending) ro variable konam
			finishedTests++
			// todo: in structure ro khosham nemiad
			// todo suspend vaghtiye ke taze sakhte shode
			// todo pending vaghtiye ke khode user suspend karde
		} else if quiz.Status == "suspend" || quiz.Status == "pending" {
			suspendedTests++
		}
		for _, answer := range quiz.UserAnswers {
			userAnswers = append(userAnswers, answer)
		}
	}
	correctAnswerCount, incorrectAnswerCount, omittedAnswerCount := M.CalculateAnswersStats(userAnswers)
	unusedQuestions := totalQuestionsCount - usedQuestionsCount
	// in some cases (no course but user has quizzes), the unused questions may be negative
	if unusedQuestions < 0 {
		unusedQuestions = 0
	}
	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"correctAnswerCount":   correctAnswerCount,
			"incorrectAnswerCount": incorrectAnswerCount,
			"omittedAnswerCount":   omittedAnswerCount,
			"createdTests":         createdTests,
			"completedTests":       finishedTests,
			"suspendedTests":       suspendedTests,
			"totalQuestionsCount":  totalQuestionsCount,
			"usedQuestionsCount":   usedQuestionsCount,
			"unusedQuestionsCount": totalQuestionsCount - usedQuestionsCount,
		},
	})
}

// report correct, incorrect and omitted answers count of every subject and system
// TODO optimize this code
func ReportQuiz(c *fiber.Ctx) error {
	// 1. get all user's quizzes
	user := c.Locals("user").(M.User)
	// options is needed in user preload in correct and incorrect answer coount
	if err := D.DB().Preload("Quizzes.UserAnswers.Question.Options").
		Preload("Quizzes.UserAnswers.Question.System.Subject").
		Find(&user).Error; err != nil {
		return U.DBError(c, err)
	}
	// 2. group answers by subject and system and calculate every answer stat
	// create this object for every system and subject
	subjects := []M.QuizAnswerStats{}
	systems := []M.QuizAnswerStats{}
	// insert every user's quizzes subject and system into subjects and systems array
	subjectIDs := []uint{}
	systemIDs := []uint{}
	for _, quiz := range user.Quizzes {
		for _, answer := range quiz.UserAnswers {
			// every answer's stat should be calculate separately for subject and system
			// if system doesn't exist already in their array create it
			if !U.ExistsInArray[uint](systemIDs, answer.Question.SystemID) {
				systemIDs = append(systemIDs, answer.Question.SystemID)
				systems = append(systems, M.QuizAnswerStats{
					ID:    answer.Question.SystemID,
					Title: answer.Question.System.Title,
				})
				quizStat := M.FindQuizAnswerStats(systems, answer.Question.SystemID)
				quizStat.UpdateQuizStat(answer)
			} else { // if this system already exists, just update it
				quizStat := M.FindQuizAnswerStats(systems, answer.Question.SystemID)
				quizStat.UpdateQuizStat(answer)
			}
			// if subject  doesn't exist already in their array create it
			if !U.ExistsInArray[uint](subjectIDs, answer.Question.System.SubjectID) {
				subjectIDs = append(subjectIDs, answer.Question.System.SubjectID)
				subjects = append(subjects, M.QuizAnswerStats{
					ID:    answer.Question.System.SubjectID,
					Title: answer.Question.System.Subject.Title,
				})
				quizStat := M.FindQuizAnswerStats(subjects, answer.Question.System.SubjectID)
				quizStat.UpdateQuizStat(answer)
			} else { // if this subject already exists, just update it
				quizStat := M.FindQuizAnswerStats(subjects, answer.Question.System.SubjectID)
				quizStat.UpdateQuizStat(answer)
			}
		}
	}
	return c.JSON(fiber.Map{"data": fiber.Map{
		"subjects": subjects,
		"systems":  systems,
	}})
}
