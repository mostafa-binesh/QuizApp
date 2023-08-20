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
func UserBoughtCoursesWithExpirationDate(userID uint) (*[]CourseWithExpirationDate, error) {
	var userBoughtCourses []CourseUser
	// get the CourseUser and preload it it course
	if err := D.DB().
		Model(&CourseUser{}).
		Where("user_id = ? AND expiration_date > ?", userID, time.Now()).
		Preload("Course.ParentCourse.Subjects.Systems").
		Find(&userBoughtCourses).
		Error; err != nil {
		return nil, err
	}
	var userCourses []CourseWithExpirationDate
	for i := 0; i < len(userBoughtCourses); i++ {
		userCourses = append(userCourses, CourseWithExpirationDate{
			ID:             userBoughtCourses[i].ID,
			Title:          userBoughtCourses[i].Course.Title,
			ExpirationDate: userBoughtCourses[i].ExpirationDate,
			Subjects:       userBoughtCourses[i].Course.ParentCourse.Subjects,
			Duration:       userBoughtCourses[i].Course.ParentCourse.Duration,
		})
	}
	return &userCourses, nil
}
