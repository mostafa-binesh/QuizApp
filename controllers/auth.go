package controllers

import (
	D "docker/database"
	M "docker/models"
	S "docker/services"
	U "docker/utils"
	"fmt"
	"strings"
	"time"

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
	if emailUniqueness := D.DB().Find(&M.User{}, "email = ?", payload.Email); emailUniqueness.RowsAffected != 0 || emailUniqueness.Error != nil {
		return U.ResErr(c, "Email already exists")
	}
	// ! here we need to check that if the order exists and then if exists
	// ! > add the courses for the user
	// get user orders
	order, err := U.WCClient().Order.Get(int64(payload.OrderID), nil)
	if err != nil {
		return U.ResErr(c, fmt.Sprint("Woocommerce Error: , ", err.Error()))
	}
	// todo: move this add user's bought courses to its service
	// add bought courses to user's courses
	var purchasedCoursesIDs []int64
	purchasedCourseIDPayDateMap := make(map[int64]time.Time)
	for _, item := range order.LineItems {
		// todo: uncomment order.Status line in production
		// if order.Status != "completed" {
		// 	continue
		// }
		purchasedCoursesIDs = append(purchasedCoursesIDs, item.ProductID)
		var payTime time.Time
		// todo: dev only, delete order.DatePaidGMT=="" on production
		// if DatePaidGMT is empty, the order hasn't been paid
		if order.DatePaidGmt == "" {
			payTime = time.Now()
		} else {
			payTime = S.ConvertWCTimeToStandard(order.DatePaidGmt)
		}
		purchasedCourseIDPayDateMap[item.ProductID] = payTime
	}
	// get the courses
	var childrenCourses []M.Course
	if result := D.DB().Find(&childrenCourses, "woocommerce_id IN ?", purchasedCoursesIDs); result.Error != nil {
		return U.DBError(c, result.Error)
	}
	// get the parent course IDs using child courses
	var parentCourseIDs []uint
	for _, childrenCourse := range childrenCourses {
		if childrenCourse.ParentID != nil {
			parentCourseIDs = append(parentCourseIDs, *childrenCourse.ParentID)
		}
	}
	// if no the order doesn't contain any valid child course, show error
	if len(childrenCourses) == 0 || len(parentCourseIDs) == 0 {
		return U.ResErr(c, "Your order doesn't contain any valid course")
	}
	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// we need to set the child courses into the course_user table
	newUser := M.User{
		Email:    payload.Email,
		Password: string(hashedPassword),
	}
	//  add created user to the database
	result := D.DB().Create(&newUser)
	if result.Error != nil {
		return U.ResErr(c, "couldn't create the user")
	}
	// create course_user records
	var userNewBoughtCourses []M.CourseUser
	for _, childCourse := range childrenCourses {
		// if the childCoruse doesn't have parent, don't insert it
		if childCourse.ParentID == nil {
			continue
		}
		userNewBoughtCourses = append(userNewBoughtCourses, M.CourseUser{
			UserID:         int(newUser.ID),
			CourseID:       int(*&childCourse.ID),
			ExpirationDate: purchasedCourseIDPayDateMap[int64(childCourse.WoocommerceID)].Add(time.Duration(childCourse.ValidityDaysPeriod) * time.Hour * 24),
		})
	}
	// we SAVE records to the database because some user may bought other courses already
	if err := D.DB().Save(&userNewBoughtCourses).Error; err != nil {
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
	fmt.Printf("payload: %v\n", payload)
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
	fmt.Printf("found courses: %v\n", courses)
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
	fmt.Printf("payload: %v\n", payload)
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
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return U.ResErr(c, loginError)
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
