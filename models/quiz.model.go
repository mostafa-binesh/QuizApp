package models

import (
	U "docker/utils"
	"strings"
	// "fmt"
	"time"
)

type QuizQuestionMode uint

const (
	AllQuestionMode QuizQuestionMode = 1 + iota
	MarkedQuestionMode
	SingleSelectQuestionMode
	MultipleSelectQuestionMode
)

type Quiz struct {
	ID     uint  `json:"id,omitempty" gorm:"primary_key"`
	UserID uint  `json:"-"`
	User   *User `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	// TODO add lesson : lesson >> ? lesson == course ?
	Status      string       `json:"status,omitempty"`
	UserAnswers []UserAnswer `json:"userAnswers,omitempty"`
	CreatedAt   time.Time    `json:"date" gorm:"not null;default:now()"`
	EndTime     *time.Time   `json:"-" gorm:"not null;default:now()"`
	Duration    uint         `json:"duration" gorm:"not null"`
	CourseID    uint         `json:"-"`
	Course      *Course      `json:"course,omitempty" gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	// mode = tutor, timed
	Mode string `json:"mode" gorm:"type:varchar(255)"`
	// type like nextGeneration
	Type string `json:"type" gorm:"type:varchar(255)"`
	// QuestionMode = All, Marked, singleSelect, multipleSelect
	QuestionMode QuizQuestionMode `json:"questionType"`
}

// used for creating new quiz
type QuizInput struct {
	QuestionsCount int              `json:"questionsCount" validate:"required,min=1"`
	SystemIDs      []uint           `json:"systemIDs" validate:"required"`
	QuizMode       []string         `json:"quizMode" validate:"required"`
	QuizType       []string         `json:"quizType" validate:"required"`
	QuestionMode   QuizQuestionMode `json:"questionMode" validate:"required"`
}

// used for listing the user's quizzes
type QuizList struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

// used to convert backend quiz model to front mocked model
type QuizToFront struct {
	ID                uint             `json:"no" gorm:"primary_key"`
	Questions         []*Question      `json:"questions"` // question with options only
	UserAnswers       [][]*string      `json:"userAnswers"`
	UserNotes         []*string        `json:"userNotes"`
	UserMarks         []bool           `json:"userMarks"`
	SubmitedQuestions []bool           `json:"submitedQuestions"`
	QuestionsStatus   []*string        `json:"questionsStatus"`
	SpentTimes        []*uint          `json:"spentTimes"`
	RemainingHours    int              `json:"remainingHours"`
	RemainingMinutes  int              `json:"remainingMinutes"`
	RemainingSeconds  int              `json:"remainingSeconds"`
	QuizState         string           `json:"quizState"`
	CreatedAt         string           `json:"date"`
	Course            string           `json:"courseName"`
	Mode              []string         `json:"mode"`
	Type              []string         `json:"type"`
	QuestionMode      QuizQuestionMode `json:"questionMode"`
	// TODO add quizID
}

// convert quiz model to mocked front quiz structure
func (quiz *Quiz) ConvertQuizToQuizToFront() QuizToFront {
	quizFront := QuizToFront{}
	quizFront.ID = quiz.ID
	var quizFrontQuestions []*Question
	var userAnswers [][]*string
	var userNotes []*string
	var userMarks []bool
	var submitedQuestions []bool
	var questionsStatus []*string
	var spentTimes []*uint
	for _, v := range quiz.UserAnswers {
		quizFrontQuestions = append(quizFrontQuestions, v.Question)
		// # handling answers
		// answers can be "A" or "A,B"
		var answersPtr []*string
		if v.Answer != nil {
			answers := strings.Split(*v.Answer, ",")
			answersPtr = U.ConvertSliceToPtrSlice[string](answers)
		}
		userAnswers = append(userAnswers, answersPtr)
		userNotes = append(userNotes, v.Note)
		userMarks = append(userMarks, v.IsMarked)
		submitedQuestions = append(submitedQuestions, v.Submitted)
		questionsStatus = append(questionsStatus, &v.Status)
		spentTimes = append(spentTimes, &v.SpentTime)
	}
	quizFront.Questions = quizFrontQuestions
	quizFront.UserAnswers = userAnswers
	quizFront.UserNotes = userNotes
	quizFront.UserMarks = userMarks
	quizFront.SubmitedQuestions = submitedQuestions
	quizFront.QuestionsStatus = questionsStatus
	quizFront.SpentTimes = spentTimes
	hour, min, sec := U.CalculateRemainingTime(quiz.Duration) // todo make this a quiz function
	quizFront.RemainingHours = hour
	quizFront.RemainingMinutes = min
	quizFront.RemainingSeconds = sec
	quizFront.QuizState = quiz.Status
	quizFront.Course = quiz.Course.Title
	quizFront.CreatedAt = quiz.CreatedAt.Format("2006-01-02T15:04:05.000Z")
	quizFront.Mode = strings.Split(quiz.Mode, ",")
	quizFront.Type = strings.Split(quiz.Type, ",")
	quizFront.QuestionMode = quiz.QuestionMode
	return quizFront
}

// convert quiz model to mocked front quiz structure
// userAnswers come from database user's quiz.userAnswers field
func (frontQuiz *QuizToFront) ConvertQuizFrontToQuiz(userAnswers []UserAnswer) []UserAnswer {
	// go through each frontQuiz, get the values and insert them into userAnswers array
	// handling user answers
	for i := range userAnswers {
		// fmt.Printf(*userAnswers[i].Answer + "\n")
		if frontQuiz.UserAnswers[i] != nil {
			// need to convert the *answers of array to answer to be able to join them using ,
			var answersString []string
			for _, answer := range frontQuiz.UserAnswers[i] {
				// because answer is a pointer, we need to check the nil pointer
				if answer != nil {
					answersString = append(answersString, *answer)
				}
			}
			joinedString := strings.Join(answersString, ",")
			userAnswers[i].Answer = &joinedString
		}
		if frontQuiz.UserNotes[i] != nil {
			userAnswers[i].Note = frontQuiz.UserNotes[i]
		}
		userAnswers[i].IsMarked = frontQuiz.UserMarks[i]
		if frontQuiz.QuestionsStatus[i] != nil {
			userAnswers[i].Status = *frontQuiz.QuestionsStatus[i]
		}
		if frontQuiz.SpentTimes[i] != nil {
			userAnswers[i].SpentTime = *frontQuiz.SpentTimes[i]
		}
		userAnswers[i].Submitted = frontQuiz.SubmitedQuestions[i]
	}
	return userAnswers
}
func (quiz *Quiz) CalculateRemainingSeconds(hours, minutes, seconds int) {
	totalSeconds := (hours * 3600) + (minutes * 60) + seconds
	quiz.Duration = uint(totalSeconds)
}

// its used for user.Courses so i needed to make the argument refrence
func ConvertQuizToQuizList(quizzes []*Quiz) []QuizList {
	var quizList []QuizList
	for i := 0; i < len(quizzes); i++ {
		quizWithList := QuizList{
			ID:     quizzes[i].ID,
			Status: quizzes[i].Status,
		}
		quizList = append(quizList, quizWithList)
	}
	return quizList
}
