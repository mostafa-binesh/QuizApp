package models

type System struct {
	ID        uint        `json:"id" gorm:"primary_key"`
	Title     string      `json:"title"`
	SubjectID uint        `json:"subjectID"`
	Subject   *Subject    `json:"subject,omitempty"`
	Questions []*Question `json:"questions,omitempty" gorm:"foreignKey:SystemID"`
}
type SystemWithQuestionsCount struct {
	ID                           uint   `json:"id" gorm:"primary_key"`
	Title                        string `json:"title"`
	SubjectID                    uint   `json:"subjectID"`
	QuestionsCount               int    `json:"questionsCount"`
	TraditionalQuestionsCount    int    `json:"traditionalQuestionsCount"`
	NextGenerationQuestionsCount int    `json:"nextGenerationQuestionsCount"`
}

func (s System) QuestionsCount() (traditionalQount, nextGenerationCount int) {
	for _, q := range s.Questions {
		if q.IsTraditional() {
			traditionalQount++
		} else {
			nextGenerationCount++
		}
	}
	return
}
