package models

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type Role uint

const (
	UserRole  Role = 1
	AdminRole Role = 2
)

// ! the model that been used for migration and retrieve and add data to the database
type User struct {
	ID          uint          `gorm:"primaryKey"`
	Email       string        `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password    string        `json:"-" gorm:"type:varchar(100);not null"`
	Role        uint          `gorm:"default:1;not null"` // 1: normal user, 2: moderator, 3: admin
	Verified    bool          `gorm:"not null;default:false"`
	CreatedAt   *time.Time    `gorm:"not null;default:now()"`
	UpdatedAt   *time.Time    `gorm:"not null;default:now()"`
	Courses     []*Course     `gorm:"many2many:course_user;"`
	Quizzes     []Quiz        `json:"quizzes" gorm:"foreignKey:UserID"`
	UserAnswers []*UserAnswer `json:"userAnswer" gorm:"foreignKey:UserID"`
}

// as the Role field is uint, we need to convert it to string sometimes
// eg. 1 returns "user"
func (c *User) RoleString() string {
	switch c.Role {
	case uint(AdminRole):
		return "admin"
	case uint(UserRole):
		return "user"
	default:
		return "none"
	}
}

type MinUser struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}
type MinUserWithCoursesIDs struct {
	ID      uint   `json:"id"`
	Email   string `json:"email"`
	Courses []uint `json:"courses"`
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
	Password   string `json:"password" validate:"omitempty,min=4"`
	CoursesIDs []uint `json:"courses" validate:"required"`
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
type AddCourseUsingOrderID struct {
	OrderID uint `json:"orderID" validate:"required"`
}

// get the authenticated user interface from fiber context locals variabels and convert to user model
// auth middleware should be done already
// didn't add it in auth utility becase of cycle import error
func AuthedUser(c *fiber.Ctx) User {
	return c.Locals("user").(User)
}
