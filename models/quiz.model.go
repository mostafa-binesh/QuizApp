package models

type Quiz struct {
	ID     uint `json:"id" gorm:"primary_key"`
	UserID uint `json:"userID"`
	User   User `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
	// TODO lesson >> ? lesson == course ?
	Status   string `json:"status"`
	CourseID uint   `json:"courseID"`
	Course   User   `json:"course" gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE"`
}

// used for creating new quiz
type QuizInput struct {
	QuestionsCount int `json:"questionsCount" validate:"required,min=1"`
}

// used for listing the user's quizzes
type QuizList struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
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
