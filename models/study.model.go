package models

import "time"

const maxStudyHour = 20
const minStudyHour = 0

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
type StudyPlanUpdateInput struct {
	ID uint `json:"id" gorm:"primaryKey" validate:"required"`
	// Date       time.Time `json:"date" validate:"required"`
	IsFinished *bool `json:"isFinished" validate:"required"`
}

// functions

// set IsFinished field to true
func (sp *StudyPlan) Finish() {
	sp.IsFinished = true
}

// check if all values are between min and max study hour
func (sp *CreateNewStudyPlanInput) ValidateWorkingHours() {
	// check if all values are between min and max study hour
	for i := range sp.WorkingHours {
		if sp.WorkingHours[i] > maxStudyHour {
			sp.WorkingHours[i] = maxStudyHour
		} else if sp.WorkingHours[i] < minStudyHour {
			sp.WorkingHours[i] = minStudyHour
		}
	}
}
