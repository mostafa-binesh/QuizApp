package admin

import (
	D "docker/database"
	F "docker/database/filters"
	M "docker/models"
	U "docker/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// ############################
// ##########    USER   #############
// ############################

// ! Index User with admin/users route
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
			ID:           user[i].ID,
			// Name:         user[i].Name,
			// PhoneNumber:  user[i].PhoneNumber,
			// PersonalCode: user[i].PersonalCode,
			// NationalCode: user[i].NationalCode,
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
	// return c.SendString("wtf")
	user := M.User{}
	payload := new(M.EditInput)
	if err := c.BodyParser(payload); err != nil {
		return U.ResErr(c, err.Error())
	}
	if errs := U.Validate(payload, c.Params("id")); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	result1 := D.DB().Where("id = ?", c.Params("id")).First(&user)
	if result1.Error != nil {
		return U.DBError(c, result1.Error)
	}
	// user.Name = payload.Name
	// user.NationalCode = payload.NationalCode
	if payload.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		if err != nil {
			return U.ResErr(c, "خطا در پردازش رمز عبور")
		}
		user.Password = string(hashedPassword)
	}
	// user.PhoneNumber = payload.PhoneNumber
	// user.PersonalCode = payload.PersonalCode
	result := D.DB().Save(&user)
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	return U.ResMessage(c, "کاربر ویرایش شد")
}

// ! user by id with admin/users/{id}
func UserByID(c *fiber.Ctx) error {
	user := M.User{}
	result := D.DB().Where("id = ?", c.Params("id")).Find(&user)
	if result.RowsAffected == 0 { // can check the same condition with user.Name == ""
		return U.ResErr(c, "کاربر وجود ندارد")
	}
	minUser := M.MinUser{
		ID:           user.ID,
		// Name:         user.Name,
		// PhoneNumber:  user.PhoneNumber,
		// PersonalCode: user.PersonalCode,
		// NationalCode: user.NationalCode,
	}
	return c.JSON(fiber.Map{
		"data": minUser,
	})
}

// ! Delete user with admin/users/{}
func DeleteUser(c *fiber.Ctx) error {
	result := D.DB().Delete(&M.User{}, c.Params("id"))
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	if result.RowsAffected == 0 {
		return U.ResErr(c, "کاربر یافت نشد")
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
	payload := new(M.SignUpInput)
	// ! parse body
	if err := c.BodyParser(payload); err != nil {
		return U.ResErr(c, err.Error())
	}
	// ! validate request
	if errs := U.Validate(payload, c.Params("id")); errs != nil {
		return U.ResValidationErr(c, errs)
	}
	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		U.ResErr(c, err.Error())
	}
	newUser := M.User{
		// Name:         payload.Name,
		// PhoneNumber:  strings.ToLower(payload.PhoneNumber),
		Password:     string(hashedPassword),
		// PersonalCode: payload.PersonalCode,
		// NationalCode: payload.NationalCode,
	}
	result := D.DB().Create(&newUser)
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	return U.ResMessage(c, "کاربر ایجاد شد") // ! TODO talk with mohsen: should i send statusCreated or 200 ?
}
