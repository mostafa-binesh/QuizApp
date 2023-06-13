package models

import (
	"time"
)

// ! the model that been used for migration and retrieve and add data to the database
type User struct {
	// ID        *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	// gorm.Model
	ID uint `gorm:"primaryKey"`
	// Name         string `gorm:"type:varchar(100);not null"`
	// PhoneNumber  string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Email    string `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password string `gorm:"type:varchar(100);not null"`
	Role     uint   `gorm:"default:1;not null"` // 1: normal user, 2: moderator, 3: admin
	// PersonalCode string `gorm:"type:varchar(10);uniqueIndex"`
	// NationalCode string `gorm:"type:varchar(10);uniqueIndex"`
	// Provider  *string    `gorm:"type:varchar(50);default:'local';not null"`
	// Photo     *string    `gorm:"not null;default:'default.png'"`
	Verified  bool       `gorm:"not null;default:false"`
	CreatedAt *time.Time `gorm:"not null;default:now()"`
	UpdatedAt *time.Time `gorm:"not null;default:now()"`
	Courses   []*Course  `gorm:"many2many:user_courses;"`
	Quizzes   []*Quiz    `json:"quizzes" gorm:"foreignKey:UserID"`
}
type MinUser struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

// ! this model has been used in signup handler
type SignUpInput struct {
	Email    string `json:"email" validate:"required,email"`
	OrderID  uint   `json:"orderId" validate:"required,numeric"`
	Password string `json:"password" validate:"required,min=4"`
}
type AdminCreateUserInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=4"`
	Courses  []uint `json:"courses" validate:"required"`
}
type AdminEditUserInput struct {
	Email    string `json:"email" validate:"required,email,dunique=users.email"`
	Password string `json:"password" validate:"required,min=4"`
	// Courses  []uint `json:"courses" validate:"required"`
}

// ! this model has been used in Edit user handler
type EditInput struct {
	Name         string `json:"name" validate:"required"`
	PhoneNumber  string `json:"phoneNumber" validate:"required,regex=^09\d{9}$,dunique=users"`
	PersonalCode string `json:"personalCode" validate:"required,max=10,numeric,dunique=users"`
	NationalCode string `json:"nationalCode" validate:"required,len=10,numeric,dunique=users"`
	Password     string `json:"password"`
	// Photo string `json:"photo"`
}

// ! this model has been used in login handler
type SignInInput struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// ! not been used
type UserResponse struct {
	ID          uint      `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	Role        string    `json:"role,omitempty"`
	Photo       string    `json:"photo,omitempty"`
	Provider    string    `json:"provider"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
