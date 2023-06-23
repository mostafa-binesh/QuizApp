package admin

import (
	D "docker/database"
	F "docker/database/filters"
	M "docker/models"
	U "docker/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// ############################
// ##########    USER   #############
// ############################

// ! Index User with admin/users route
// ! admin/users
func IndexUser(c *fiber.Ctx) error {
	user := []M.User{}
	pagination := U.ParsedPagination(c)
	D.DB().Scopes(
		F.FilterByType(
			F.FilterType{QueryName: "name", Operator: "LIKE"},
			F.FilterType{QueryName: "nationalCode", ColumnName: "national_code"},
			F.FilterType{QueryName: "personalCode", ColumnName: "personal_code"}),
		U.Paginate(user, pagination)).Find(&user)
	pass_data := []M.MinUser{}
	for i := 0; i < len(user); i++ {
		pass_data = append(pass_data, M.MinUser{
			ID:    user[i].ID,
			Email: user[i].Email,
		})
	}
	return c.JSON(fiber.Map{
		"meta": pagination,
		"data": pass_data,
	}) // same as return U.ResWithPagination(c, pass_data, *pagination)
}
func CheckPasswordHash(password string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
func EditUser(c *fiber.Ctx) error {
	payload := new(M.AdminEditUserInput)
	// parse body
	if err := c.BodyParser(payload); err != nil {
		return U.ResErr(c, err.Error())
	}
	// validate the payload
	if errs := U.Validate(payload, c.Params("id")); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	// get the user with given id in param
	user := M.User{}
	if err := D.DB().Where("id = ?", c.Params("id")).First(&user).Error; err != nil {
		return U.DBError(c, err)
	}
	// hash the password if password exist in payload
	if payload.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		if err != nil {
			return U.ResErr(c, "خطا در پردازش رمز عبور")
		}
		user.Password = string(hashedPassword)
	}
	// create or update user's courses:
	// Get the courses by IDs
	var courses []M.Course
	D.DB().Where("id IN (?)", payload.CoursesIDs).Find(&courses)
	// Update the user's courses
	D.DB().Model(&user).Association("Courses").Replace(&courses)
	// save the user
	result := D.DB().Save(&user)
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	return U.ResMessage(c, "User updated successfully")
}

// ! user by id with admin/users/{email}
func UserByEmail(c *fiber.Ctx) error {
	user := M.User{}
	result := D.DB().Where("email = ?", c.Params("email")).Preload("Courses").First(&user)
	if result.RowsAffected == 0 { // can check the same condition with user.Name == ""
		return U.ResErr(c, "User doesn't exist")
	}
	// extract user's courses ids
	var userCoursesIDs []uint
	for _, course := range user.Courses {
		userCoursesIDs = append(userCoursesIDs, course.ID)
	}
	minUserWithCoursesIDs := M.MinUserWithCoursesIDs{
		ID:      user.ID,
		Email:   user.Email,
		Courses: userCoursesIDs,
	}
	return c.JSON(fiber.Map{
		"data": minUserWithCoursesIDs,
	})
}

// ! Delete user with admin/users/{}
func DeleteUser(c *fiber.Ctx) error {
	user := new(M.User)
	result := D.DB().Preload("Courses").First(user, c.Params("id"))
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	if result.RowsAffected == 0 {
		return U.ResErr(c, "کاربر یافت نشد")
	}
	// Delete all courses associated with the user
	// err := D.DB().Delete(&user.Courses).Error // deletes the courses
	err := D.DB().Model(&user).Association("Courses").Clear() // deletes the user_courses
	if err != nil {
		return U.DBError(c, err)
	}
	if err = D.DB().Delete(&M.User{}, c.Params("id")).Error; err != nil {
		return U.DBError(c, err)
	}
	return U.ResMessage(c, "کاربر حذف شد")
}

func UserVerification(c *fiber.Ctx) error {
	result := D.DB().Model(&M.User{}).Where("id = ?", c.Params("id")).Update("verified", true)
	if result.Error != nil {
		U.DBError(c, result.Error)
	}
	return U.ResMessage(c, "کاربر تایید شد")
}
func UserUnVerification(c *fiber.Ctx) error {
	result := D.DB().Model(&M.User{}).Where("id = ?", c.Params("id")).Update("verified", false)
	if result.Error != nil {
		U.DBError(c, result.Error)
	}
	return U.ResMessage(c, "کاربر رد شد")
}
func AddUser(c *fiber.Ctx) error {
	payload := new(M.AdminCreateUserInput)
	// parsing the payload
	if err := c.BodyParser(payload); err != nil {
		U.ResErr(c, err.Error())
	}
	// validation the payload
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	fmt.Printf("payload: %v\n", payload)
	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return U.ResErr(c, err.Error())
	}
	var courses []*M.Course
	fmt.Printf("till here")
	if err := D.DB().Find(&courses, payload.Courses).Error; err != nil {
		return U.DBError(c, err)
	}
	newUser := M.User{
		Email:    payload.Email,
		Password: string(hashedPassword),
		// Courses:  courses,
	}
	result := D.DB().Create(&newUser)
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	// add courses to the user
	D.DB().Model(&newUser).Association("Courses").Replace(&courses)
	return U.ResMessage(c, "User has been created")
}
