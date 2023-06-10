package models

type Course struct {
	ID            uint   `json:"id" gorm:"primary_key"`
	WoocommerceID uint   `json:"woocommerce_id"`
	Title         string `json:"title"`
	Users  []*User `gorm:"many2many:user_courses;"`
}

// model used for creating new course
type CourseInput struct {
	ID            uint   `json:"id" validate:"required"`
	WoocommerceID uint   `json:"wc_id" validate:"required"`
	Title         string `json:"title" validate:"required"`
}
