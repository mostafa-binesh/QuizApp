package models

import (
	"time"

	"gorm.io/gorm"
)

type CourseUser struct {
	UserID         int `gorm:"foreignKey:ID"`
	CourseID       int `gorm:"foreignKey:ID"`
	ExpirationDate time.Time
}

// Set the table name for the join table
func (CourseUser) TableName() string {
	return "course_user"
}

// Custom global scope for UserCourse model to filter out expired records
func NonExpiredCourses(db *gorm.DB) *gorm.DB {
	currentTime := time.Now()
	return db.Where("expiration_date > ?", currentTime)
}
