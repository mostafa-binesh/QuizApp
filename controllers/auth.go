package controllers

import (
	"fmt"
	// "github.com/go-playground/validator/v10"
	D "docker/database"
	M "docker/models"
	U "docker/utils"

	"github.com/gofiber/fiber/v2"

	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Person struct {
	Name string `json:"name" xml:"name" form:"name"`
	Pass string `json:"pass" xml:"pass" form:"pass"`
}

func SignUpUser(c *fiber.Ctx) error {
	payload := new(M.SignUpInput)
	// parsing the payload
	if err := c.BodyParser(payload); err != nil {
		U.ResErr(c, err.Error())
	}
	fmt.Printf("payload: %v\n", payload)
	// ! validate request
	// ! todo: create a shorter function for the validation, like payload.validate(), validate function can get a T template
	// validate the payload
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	// ! here we need to check that if the order exists and then if exists
	// ! > add the courses for ther user
	order, err := U.WCClient().Order.Get(int64(payload.OrderID), nil)
	if err != nil {
		return U.ResErr(c, fmt.Sprint("خطای ووکامرس, ", err.Error()))
	}
	// var courses []M.Course
	var purchasedOrderIDs []int64
	for _, item := range order.LineItems {
		// fmt.Println(item.ProductID)
		purchasedOrderIDs = append(purchasedOrderIDs, item.ProductID)
		// courses = append(courses, M.Course{
		// 	WoocommerceID: uint(item.ProductID),
		// })
	}
	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// get the courses
	var courses []*M.Course
	if result := D.DB().Find(&courses, "woocommerce_id IN ?", purchasedOrderIDs); result.Error != nil {
		return U.DBError(c, result.Error)
	}
	fmt.Printf("found courses: %v\n", courses)
	// inja miad va course ha ham mikhad besze
	// vali faghay bayad append kone
	newUser := M.User{
		// Name:         payload.Name,
		Email: payload.Email,
		// PhoneNumber:  strings.ToLower(payload.PhoneNumber), // ! can use fiber toLower function that has better performance
		Password: string(hashedPassword),
		// PersonalCode: payload.PersonalCode,
		// NationalCode: payload.NationalCode,
		// Photo:    &payload.Photo, // ? don't know why add & in the payload for photo
		Courses: courses,
	}
	// ! add user to the database
	result := D.DB().Create(&newUser)
	// ! if any error exist in the create process, write the error
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "couldn't create the user"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"msg": "user has been created successfully"})

}

func Login(c *fiber.Ctx) error {
	payload := new(M.SignInInput)
	// ! parse payload
	if err := c.BodyParser(payload); err != nil {
		U.ResErr(c, err.Error())
	}
	// ! validate request
	if errs := U.Validate(payload); errs != nil {
		return U.ResValidationErr(c, errs)
	}
	var user M.User
	result := D.DB().First(&user, "email = ?", strings.ToLower(payload.Email))
	// ! the reason we didn't handle the error first,
	// ! - is because not found return error option is disabled
	if result.RowsAffected == 0 {
		// return ReturnError(c, "ایمیل یا رمز عبور اشتباه است")
		return U.ResErr(c, "کد پرسنلی یا رمز عبور اشتباه است")
	}
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	// ! compare the password of payload and returned user from database
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return U.ResErr(c, "کد پرسنلی یا رمز عبور اشتباه است")
	}
	sess := U.Session(c)
	sess.Set(U.USER_ID, user.ID)
	if err := sess.Save(); err != nil {
		return U.ResErr(c, "خطا در ورود")
	}
	return U.ResMessage(c, "ورود انجام شد")
}
func Logout(c *fiber.Ctx) error {
	// ! just removing the session
	sess, err := U.Store.Get(c)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "not authenticated",
		})
	}
	if err := sess.Destroy(); err != nil {
		panic(err)
	}
	return c.SendString("logged out successfully")
}
func Dashboard(c *fiber.Ctx) error {
	// ! has a AuthMiddleware before here
	// ! if session and user exists, client can access here
	user := c.Locals("user").(M.User)
	return c.JSON(fiber.Map{"dashboard": "heres the dashboard", "user": user})
}
func AuthMiddleware(c *fiber.Ctx) error {
	sess := U.Session(c)
	userID := sess.Get(U.USER_ID)
	if userID == nil {
		return ReturnError(c, "not authenticated", fiber.StatusUnauthorized) // ! notAuthorized is notAuthenticated
	} else {
		c.SendString(fmt.Sprintf("user id is: %s", sess.Get(U.USER_ID)))
	}
	var user M.User
	result := D.DB().Find(&user, userID)
	if result.Error != nil {
		err := sess.Destroy()
		if err != nil {
			panic(err)
		}
		return ReturnError(c, "cannot authenticate. session removed", 500)
	}
	c.Locals("user", user)
	return c.Next()

}
