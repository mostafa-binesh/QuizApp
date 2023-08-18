package models

import "time"

type StudyPlan struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Date       time.Time `json:"date"`
	Hours      uint      `json:"hours"`
	IsFinished bool      `json:"isFinished"`
	// # REFERENCES
	UserID uint `json:"-"`
	User   User `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
}

type CreateNewStudyPlanInput struct {
	StartDate    time.Time `json:"startDate" validate:"required"`
	EndDate      time.Time `json:"endDate" validate:"required"`
	WorkingHours []uint    `json:"workingHours" validate:"required"`
}

type VerifyStudyPlanDateInput struct {
	Date []time.Time `json:"date" validate:"required"`
}

// functions

// set IsFinished field to true
func (sp *StudyPlan) Finish() {
	sp.IsFinished = true
}
