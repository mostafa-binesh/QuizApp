package models

import "time"

type Guest struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Email     string     `json:"email" validate:"required"`
	FullName  string     `json:"fullName" validate:"required"`
	CreatedAt *time.Time `json:"createdAt" gorm:"not null;default:now()"`
	UpdatedAt *time.Time `json:"updatedAt" gorm:"not null;default:now()"`
}
