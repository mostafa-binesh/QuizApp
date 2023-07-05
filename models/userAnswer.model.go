package models

type UserAnswer struct {
	ID         uint      `json:"id,omitempty" gorm:"primary_key"`
	QuestionID uint      `json:"-"`
	Question   *Question `json:"question,omitempty" gorm:"foreignKey:QuestionID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	Note       *string   `json:"note,omitempty"`
	IsMarked   bool      `json:"isMarked,omitempty" gorm:"default:false;"`
	Submitted  bool      `json:"submitted,omitempty" gorm:"default:false;"`
	Status     string    `json:"status,omitempty"`
	SpentTime  uint      `json:"spentTime,omitempty"`
	UserID     uint      `json:"-"`
	User       *User     `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	QuizID     uint      `json:"-"`
	Quiz       *Quiz     `json:"quiz,omitempty" gorm:"foreignKey:QuizID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	// ! multiple-choice option
	// Answers    []*Answer `json:"answers,omitempty"`
	// ! single-choice option
	Answer *string `json:"answer,omitempty"`
}
type AnswerNote struct {
	ID       uint      `json:"id"` // is parent UserAnswer's ID
	Question *Question `json:"question"`
	QuizID   uint      `json:"quidID"`
	Note     *string   `json:"note"`
}
type EditNoteInput struct {
	Note *string `json:"note" validator:"required"`
}

// checks answer is true, false or empty
func CalculateAnswerStats(answer UserAnswer) (correctAnswerCount, incorrectAnswerCount, omittedAnswerCount uint) {
	// if answer is empty, increase omittedAnswerCount by one and return stats
	if answer.Answer == nil {
		omittedAnswerCount++
		return
	}
	// if answer is not empty, find out that answer is correct or not
	for _, option := range answer.Question.Options {
		optionIndex := option.Index
		userAnswerIndex := answer.Answer
		if *userAnswerIndex == optionIndex {
			correctAnswerCount++
			return
		}
	}
	// if answer was not correct, increase the incorrectAnswersCount be one and return
	incorrectAnswerCount++
	return
}

// array version of CalculateAnswerStats
// calculate correct, incorrect and omitted answers count of userAnswers
func CalculateAnswersStats(answers []UserAnswer) (correctAnswerCount, incorrectAnswerCount, omittedAnswerCount uint) {
	var found bool
	for _, answer := range answers {
		found = false
		if answer.Answer == nil {
			omittedAnswerCount++
			continue
		}
		for _, option := range answer.Question.Options {
			// had to make option a separated value to be able to compare with answer.answer
			optionIndex := option.Index
			userAnswerIndex := answer.Answer
			if *userAnswerIndex == optionIndex {
				correctAnswerCount++
				found = true
				break
			}
		}
		if !found {
			incorrectAnswerCount++
		}
	}
	return
}

// this model has been used in overall report handler
type QuizAnswerStats struct {
	ID                   uint   `json:"id"`
	Title                string `json:"title"`
	CorrectAnswerCount   uint   `json:"correctAnswerCount"`
	InCorrectAnswerCount uint   `json:"incorrectAnswerCount"`
	OmittedAnswerCount   uint   `json:"omittedAnswerCount"`
}

// find QuizAnswerStat with id of ID
func FindQuizAnswerStats(subjects []QuizAnswerStats, ID uint) *QuizAnswerStats {
	for i := 0; i < len(subjects); i++ {
		if subjects[i].ID == ID {
			return &subjects[i]
		}
	}
	return nil
}

// update stats of a quizAnswerStat with given correct, incorrect and omitted answers count
// func (quizStat *QuizAnswerStats) UpdateQuizStat(correctAnswerCount, incorrectAnswerCount, omittedAnswerCount uint) {
func (quizStat *QuizAnswerStats) UpdateQuizStat(answer UserAnswer) {
	correctAnswerCount, incorrectAnswerCount, omittedAnswerCount := CalculateAnswerStats(answer)
	quizStat.CorrectAnswerCount += correctAnswerCount
	quizStat.InCorrectAnswerCount += incorrectAnswerCount
	quizStat.OmittedAnswerCount += omittedAnswerCount
}
