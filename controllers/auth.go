package controllers

import (
	D "docker/database"
	M "docker/models"
	S "docker/services"
	U "docker/utils"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
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
	// validate the payload
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	// check email uniqueness
	if err := S.CheckEmailUniqueness(payload.Email); err != nil {
		return U.ResErr(c, err.Error())
	}
	// hash the password
	hashedPassword, err := S.GenerateHashedPassword(payload.Password)
	if err != nil {
		return U.ResErr(c, err.Error())
	}
	// get all courses that user have bought using orderID
	childCourses, purchasedCourseIDPayDateMap, err := S.ImportUserCoursesUsingOrderID(payload.OrderID)
	if err != nil {
		return U.ResDebug(c, err, "Failed to retreieve bought courses from woocommerce")
	}
	// create new user
	newUser := M.User{
		Email:    payload.Email,
		Password: string(hashedPassword),
		// Courses:  parentCourses,
	}
	//  add created user to the database
	result := D.DB().Create(&newUser)
	if result.Error != nil {
		return U.ResErr(c, "couldn't create the user")
	}
	// ! here we need to check that if the order exists and then if exists
	// ! > add the courses for the user
	// get the bought courses in course_user format
	CourseUser := S.AddCourseUserUsingCourses(childCourses, purchasedCourseIDPayDateMap, newUser.ID)
	// we SAVE records to the database because some user may bought other courses already
	if err := D.DB().Save(&CourseUser).Error; err != nil {
		// delete created user
		D.DB().Delete(&newUser)
		return U.DBError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"msg": "User has been created successfully"})
}
func DevsSignUpUser(c *fiber.Ctx) error {
	payload := new(M.SignUpInput)
	// parsing the payload
	if err := c.BodyParser(payload); err != nil {
		U.ResErr(c, err.Error())
	}
	// validate the payload
	// ! here we need to check that if the order exists and then if exists
	// ! > add the courses for the user
	// get user orders
	order, err := U.WCClient().Order.Get(int64(payload.OrderID), nil)
	if err != nil {
		return U.ResErr(c, fmt.Sprint("Woocommerce Error: , ", err.Error()))
	}
	// add bought courses to user's courses
	var courses []*M.Course
	var purchasedOrderIDs []int64
	for _, item := range order.LineItems {
		purchasedOrderIDs = append(purchasedOrderIDs, item.ProductID)
		courses = append(courses, &M.Course{
			WoocommerceID: uint(item.ProductID),
		})
	}
	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// get the courses
	if result := D.DB().Find(&courses, "woocommerce_id IN ?", purchasedOrderIDs); result.Error != nil {
		return U.DBError(c, result.Error)
	}
	newUser := M.User{
		Email:    payload.Email,
		Password: string(hashedPassword),
		Courses:  courses,
	}
	//  add created user to the database
	result := D.DB().Create(&newUser)
	//  if any error exist in the create process, write the error
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "couldn't create the user"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"msg": "user has been created successfully"})

}
func Devs2SignUpUser(c *fiber.Ctx) error {
	payload := new(M.SignUpInput)
	// parsing the payload
	if err := c.BodyParser(payload); err != nil {
		U.ResErr(c, err.Error())
	}
	// validate the payload
	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	newUser := M.User{
		Email:    payload.Email,
		Password: string(hashedPassword),
	}
	//  add created user to the database
	result := D.DB().Create(&newUser)
	//  if any error exist in the create process, write the error
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "couldn't create the user"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"msg": "user has been created successfully"})

}

var loginError string = "Invalid email or password"

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
		return U.ResErr(c, loginError)
	}
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	// ! compare the password of payload and returned user from database
	if err := S.CheckHashedPassword(user.Password, payload.Password); err != nil {
		return U.ResErr(c, err.Error())
	}
	// login successful, setting up session
	sess := U.Session(c)
	sess.Set(U.USER_ID, user.ID)
	if err := sess.Save(); err != nil {
		return U.ResErr(c, "Login error")
	}
	return c.JSON(fiber.Map{"data": fiber.Map{"role": user.RoleString()}})
}
func Logout(c *fiber.Ctx) error {
	// ! just removing the session
	sess, err := U.Store.Get(c)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Not Authenticated",
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
	user := M.AuthedUser(c)
	return c.JSON(fiber.Map{"dashboard": "heres the dashboard", "user": user})
}
func AuthMiddleware(c *fiber.Ctx) error {
	// get the session using session utility
	sess := U.Session(c)
	// retreive userID from session
	userID := sess.Get(U.USER_ID)
	// if userID doesn't exist, show err
	if userID == nil {
		return ReturnError(c, "not authenticated", fiber.StatusUnauthorized) // ! notAuthorized is notAuthenticated
	}
	// get user with userID from database
	var user M.User
	result := D.DB().First(&user, userID)
	if result.Error != nil || result.RowsAffected == 0 {
		// if user doesn't exist or any error happend, remove the session and show error
		// sess.Destroy() returns an error, but we don't need it here i guess
		sess.Destroy()
		return U.ResErr(c, "cannot authenticate. Please login again", fiber.StatusInternalServerError)
	}
	// if everything was ok, save the db user to locals variable "user"
	c.Locals("user", user)
	return c.Next()

}

// roleCheck middleware
// authentication must be done already
func RoleCheck(roles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(M.User)
		// check if authentication is done already, just in case
		if !ok {
			return U.ResErr(c, "Please login first", fiber.StatusUnauthorized)
		}
		// get string of user's role
		userRoleString := user.RoleString()
		// check roles
		for _, role := range roles {
			if role == userRoleString {
				// proceed to next handler
				return c.Next()
			}
		}
		// if roles don't match, show err
		return U.ResErr(c, "Access Forbidden", fiber.StatusForbidden)
	}
}
