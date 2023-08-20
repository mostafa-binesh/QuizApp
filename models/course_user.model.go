package models

import (
	"time"

	"gorm.io/gorm"
)

type CourseUser struct {
	ID             uint   `json:"id" gorm:"primary_key"`
	UserID         int    `gorm:"foreignKey:ID"`
	CourseID       int    `gorm:"foreignKey:ID"`
	Course         Course `json:"course"`
	ExpirationDate time.Time
}

// Set the table name
// if don't set it, it will set "course_users" name for the table
func (CourseUser) TableName() string {
	return "course_user"
}

// Custom global scope for CourseUser model to filter out expired records
func UserNonExpiredCourses() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		currentTime := time.Now()
		return db.Where("expiration_date > ?", currentTime)
	}
}
