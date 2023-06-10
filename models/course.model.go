package models

type Course struct {
	ID            uint    `json:"id" gorm:"primary_key"`
	WoocommerceID uint    `json:"woocommerce_id"`
	Title         string  `json:"title"`
	Users         []*User `gorm:"many2many:user_courses;"`
}

// model used for creating new course
type CourseInput struct {
	ID            uint   `json:"id" validate:"required"`
	WoocommerceID uint   `json:"wc_id" validate:"required"`
	Title         string `json:"title" validate:"required"`
}
type CourseWithTitleOnly struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Title string `json:"title"`
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
