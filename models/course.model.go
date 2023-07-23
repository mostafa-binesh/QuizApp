package models

type Course struct {
	ID            uint       `json:"id" gorm:"primary_key"`
	WoocommerceID uint       `json:"woocommerce_id" gorm:"unique"`
	Title         string     `json:"title"`
	Users         []*User    `json:"-" gorm:"many2many:course_user;"`
	Subjects      []*Subject `json:"subjects" gorm:"foreignKey:CourseID"`
	Duration      uint64     `json:"-"` // todo don't show it for now, fix it later
	ParentID      uint
	ParentCourse  *Course `gorm:"references:ParentID"` // use Company.CompanyID as references
}

// model used for creating new course
type CourseInput struct {
	WoocommerceID uint   `json:"wc_id" validate:"required"`
	Title         string `json:"title" validate:"required"`
}
type CourseWithTitleOnly struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Title string `json:"title"`
}
type CourseWithQuestionsCount struct {
	ID             uint                         `json:"id" gorm:"primary_key"`
	Title          string                       `json:"title"`
	Subjects       []*SubjectWithQuestionsCount `json:"subjects" gorm:"foreignKey:CourseID"`
	Duration       uint64                       `json:"-"` // todo don't show it for now, fix it later
	QuestionsCount int                          `json:"questionsCount"`
}

// its used for user.Courses so i needed to make the argument refrence
func ConvertCourseToCourseWithTitleOnly(courses []*Course) []CourseWithTitleOnly {
	var coursesWithTitleOnly []CourseWithTitleOnly
	for i := 0; i < len(courses); i++ {
		courseWithTitleOnly := CourseWithTitleOnly{
			ID:    courses[i].ID,
			Title: courses[i].Title,
		}
		coursesWithTitleOnly = append(coursesWithTitleOnly, courseWithTitleOnly)

	}
	return coursesWithTitleOnly
}
func ConvertCourseToCourseWithQuestionsCounts(courses []*Course) (coursesWithQuestionsCount []CourseWithQuestionsCount) {
	// because we're getting the *courses in making the array, we need to check the null pointer
	coursesWithQuestionsCount = make([]CourseWithQuestionsCount, len(courses))
	// var courseQuestionsCount int
	for _, course := range courses {
		var subjectsQuestionsCount int
		var systemsQuestionsCount int
		var tempSystems []*SystemWithQuestionsCount
		var tempSubjects []*SubjectWithQuestionsCount
		for _, subject := range course.Subjects {
			for _, system := range subject.Systems {
				systemsQuestionsCount += len(system.Questions)
				tempSystems = append(tempSystems, &SystemWithQuestionsCount{
					ID:             system.ID,
					Title:          system.Title,
					SubjectID:      system.SubjectID,
					QuestionsCount: len(system.Questions),
				})
			}
			subjectsQuestionsCount += systemsQuestionsCount
			tempSubjects = append(tempSubjects, &SubjectWithQuestionsCount{
				ID:             subject.ID,
				Title:          subject.Title,
				Systems:        tempSystems,
				CourseID:       subject.CourseID,
				QuestionsCount: systemsQuestionsCount,
			})
		}
		coursesWithQuestionsCount = append(coursesWithQuestionsCount, CourseWithQuestionsCount{
			ID:             course.ID,
			Title:          course.Title,
			Subjects:       tempSubjects,
			Duration:       0, // todo: hardcoded
			QuestionsCount: subjectsQuestionsCount,
		})
	}
	return
}
