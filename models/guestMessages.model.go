package models

import "time"

type GuestMessage struct {
	ID int `json:"id" gorm:"primaryKey"`
	// GuestID     int `json:"guestID" gorm:"uniqueIndex"`
	GuestChatID int `json:"guestChatID" gorm:"index"`
	// Guest     Guest      `json:"guest" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Sender    int8       `json:"sender"` // 1: guest, 2: admin
	Body      string     `json:"body"`
	CreatedAt *time.Time `json:"createdAt" gorm:"not null;default:now()"`
}

// payload
type GuestMessageCreate struct {
	Body string `json:"body" validate:"required,max=255"`
}
