package models

type Course struct {
	ID            uint       `json:"id" gorm:"primary_key"`
	WoocommerceID uint       `json:"woocommerce_id" gorm:"unique"`
	Title         string     `json:"title"`
	Users         []*User    `json:"-" gorm:"many2many:user_courses;"`
	Subjects      []*Subject `json:"subjects" gorm:"foreignKey:CourseID"`
	Duration      uint64     `json:"-"` // todo don't show it for now, fix it later
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
