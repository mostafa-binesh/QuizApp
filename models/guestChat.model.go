package models

import "time"

var GuestChatCategories = map[uint16]string{
	1: "A",
	2: "B",
	3: "C",
	4: "D",
	5: "E",
	6: "F",
	7: "G",
}

type GuestChat struct {
	ID        int            `json:"id" gorm:"primaryKey"`
	Category  int            `json:"category" validate:"required"` // ! needs to be in GuestChatCategories. how to add it?
	Title     string         `json:"title" validate:"required,max=255"`
	Messages  []GuestMessage `json:"msgs"`
	CreatedAt *time.Time     `json:"createdAt" gorm:"not null;default:now()"`
}

// payload
type GuestChatCreate struct {
	Title    string `json:"title" validate:"required,max=255"`
	Body     string `json:"body" validate:"required,max=255"`         // ! to create the first message
	Category int    `json:"category" validate:"required,min=1,max=7"` // ! needs to be in GuestChatCategories. how to add it?
}
