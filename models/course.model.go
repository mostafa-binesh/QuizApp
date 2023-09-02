package models

import (
	D "docker/database"
	U "docker/utils"
	"time"
)

type Course struct {
	ID                 uint       `json:"id" gorm:"primary_key"`
	WoocommerceID      uint       `json:"woocommerce_id" gorm:"uniqueIndex"`
	Title              string     `json:"title"`
	Users              []*User    `json:"-" gorm:"many2many:course_user;"`
	Subjects           []*Subject `json:"subjects" gorm:"foreignKey:CourseID"`
	Duration           uint64     `json:"duration"`
	ParentID           *uint
	ParentCourse       *Course `json:"-" gorm:"foreignKey:ParentID"` // use Company.CompanyID as references
	ValidityDaysPeriod uint    `json:"-"`
}

// used in user.courses route
type CourseWithExpirationDate struct {
	ID             uint       `json:"id"`
	Title          string     `json:"title"`
	Subjects       []*Subject `json:"subjects"`
	ExpirationDate time.Time  `json:"expirationDate"`
	Duration       uint64     `json:"duration"`
}
type CourseWithExpirationDateAndQuestionsCount struct {
	ID             uint                         `json:"id"`
	Title          string                       `json:"title"`
	Subjects       []*SubjectWithQuestionsCount `json:"subjects"`
	ExpirationDate time.Time                    `json:"expirationDate"`
	Duration       uint64                       `json:"duration"`
	QuestionsCount int                          `json:"questionsCount"`
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
	ID             uint                         `json:"id"`
	Title          string                       `json:"title"`
	Subjects       []*SubjectWithQuestionsCount `json:"subjects"`
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

// input course must be parentCourse
// ! subjects, systems and systems.questions must be preloaded
func ConvertCourseToCourseWithQuestionsCounts(course Course) (coursesWithQuestionsCount CourseWithQuestionsCount) {
	// because we're getting the *courses in making the array, we need to check the null pointer
	// var courseQuestionsCount int
	var subjectsQuestionsCount int
	var systemsQuestionsCount int
	var systemWithQuestionsCount []*SystemWithQuestionsCount
	var subjectsWithQuestionsCount []*SubjectWithQuestionsCount
	for _, subject := range course.Subjects {
		systemWithQuestionsCount = nil
		systemsQuestionsCount = 0
		for _, system := range subject.Systems {
			systemsQuestionsCount += len(system.Questions)
			systemWithQuestionsCount = append(systemWithQuestionsCount, &SystemWithQuestionsCount{
				ID:             system.ID,
				Title:          system.Title,
				SubjectID:      system.SubjectID,
				QuestionsCount: len(system.Questions),
			})
		}
		subjectsQuestionsCount += systemsQuestionsCount
		subjectsWithQuestionsCount = append(subjectsWithQuestionsCount, &SubjectWithQuestionsCount{
			ID:             subject.ID,
			Title:          subject.Title,
			Systems:        systemWithQuestionsCount,
			CourseID:       subject.CourseID,
			QuestionsCount: systemsQuestionsCount,
		})
	}
	return CourseWithQuestionsCount{
		ID:             course.ID,
		Title:          course.Title,
		Subjects:       subjectsWithQuestionsCount,
		Duration:       course.Duration,
		QuestionsCount: systemsQuestionsCount,
	}
}

// get all user's bought courses id from course_user table which are not expired
func RetrieveUserBoughtCoursesIDs(userID uint) ([]uint, error) {
	var courseIDs []uint
	if err := D.DB().Model(&CourseUser{}).Where("user_id = ? AND expiration_date > ?", userID, time.Now()).
		Pluck("course_id", &courseIDs).Error; err != nil {
		return nil, err
	}
	return courseIDs, nil
}

// get all user's bought parent courses id from course_user table which are not expired
func RetrieveUserBoughtParentCoursesIDs(userID uint) ([]uint, error) {
	var courseIDs []uint
	if err := D.DB().Model(&CourseUser{}).Where("user_id = ? AND expiration_date > ?", userID, time.Now()).
		Pluck("course_id", &courseIDs).Error; err != nil {
		return nil, err
	}
	// find all non-parent courses
	if err := D.DB().Model(&Course{}).Where("id IN ? AND parent_id IS NOT NULL", courseIDs).
		Pluck("parent_id", &courseIDs).Error; err != nil {
		return nil, err
	}
	return courseIDs, nil
}

// first gets user's bought courses from course_user table using RetrieveUserBoughtCoursesIDs function
// then get the bought courses from courses table
func RetrieveUserBoughtCourses(userID uint) ([]Course, error) {
	courseIDs, err := RetrieveUserBoughtCoursesIDs(userID)
	if err != nil {
		return nil, err
	}
	// find all courses where their id is in courseIDs
	var userBoughtCourses []Course
	if err := D.DB().Model(&Course{}).
		Find(&userBoughtCourses, "id IN ?", courseIDs).Error; err != nil {
		return nil, err
	}
	return userBoughtCourses, nil
}
func UserHasCourse(userID uint, courseID uint) (bool, error) {
	courseIDs, err := RetrieveUserBoughtCoursesIDs(userID)
	if err != nil {
		return false, err
	}
	return U.ExistsInArray[uint](courseIDs, courseID), nil
}
func UserBoughtCoursesWithExpirationDateAndQuestionsCount(userID uint) (*[]CourseWithExpirationDateAndQuestionsCount, error) {
	var userBoughtCourses []CourseUser
	// get the CourseUser and preload its course
	if err := D.DB().
		Model(&CourseUser{}).
		Where("user_id = ? AND expiration_date > ?", userID, time.Now()).
		Preload("Course.ParentCourse.Subjects.Systems.Questions").
		Find(&userBoughtCourses).
		Error; err != nil {
		return nil, err
	}
	var userCourses []CourseWithExpirationDateAndQuestionsCount
	// loop through each CourseUser model
	for i := 0; i < len(userBoughtCourses); i++ {
		courseWithQuestionsCount := ConvertCourseToCourseWithQuestionsCounts(*userBoughtCourses[i].Course.ParentCourse)
		userCourses = append(userCourses, CourseWithExpirationDateAndQuestionsCount{
			ID:             userBoughtCourses[i].ID,
			Title:          userBoughtCourses[i].Course.Title,
			ExpirationDate: userBoughtCourses[i].ExpirationDate, // most important part
			Subjects:       courseWithQuestionsCount.Subjects,   // subject is not from Subject model
			Duration:       userBoughtCourses[i].Course.ParentCourse.Duration,
			QuestionsCount: courseWithQuestionsCount.QuestionsCount,
		})
	}
	return &userCourses, nil
}
