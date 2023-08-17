package models

import "time"

type StudyPlan struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Date       time.Time `json:"date"`
	Hours      uint      `json:"hours"`
	IsFinished bool      `json:"isFinished"`
	// # REFERENCES
	UserID uint `json:"-"`
	User   User `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
}
