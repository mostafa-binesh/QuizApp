package controllers

import (
	D "docker/database"
	// "math/rand"
	// "time"

	// F "docker/database/filters"
	M "docker/models"
	U "docker/utils"
	"fmt"

	WC "github.com/chenyangguang/woocommerce"
	"github.com/gofiber/fiber/v2"
)

// todo: with counts ro az comment dar biar
func AllCourses(c *fiber.Ctx) error {
	user := c.Locals("user").(M.User)
	// find all bought courses ids
	userBoughtCourses, err := M.UserBoughtCoursesWithExpirationDate(user.ID)
	if err != nil {
		return U.DBError(c, err)
	}
	return c.JSON(fiber.Map{"data": userBoughtCourses})
}

// all subject of course with id of courseID
func CourseSubjects(c *fiber.Ctx) error {
	// get authenticated user
	user := c.Locals("user").(M.User)
	// get user's course with id of param with subject with system of the user
	result := D.DB().Preload("Courses", "id = ?", c.Params("courseID")).Preload("Courses.Subjects.Systems").Find(&user)
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	if len(user.Courses) != 1 {
		return U.ResErr(c, "Course not found")
	}
	// change subject model to subjectWithSystems model
	subjectWithSystems := []M.SubjectWithSystems{}
	for i := 0; i < len(user.Courses); i++ {
		for j := 0; j < len(user.Courses[i].Subjects); j++ {
			subjectWithSystems = append(subjectWithSystems, M.SubjectWithSystems{
				ID:       user.Courses[i].Subjects[j].ID,
				Title:    user.Courses[i].Subjects[j].Title,
				Systems:  user.Courses[i].Subjects[j].Systems,
				CourseID: user.Courses[i].Subjects[j].CourseID,
			})
		}
	}
	return c.JSON(fiber.Map{"data": subjectWithSystems})
}

func UpdateUserCourses(c *fiber.Ctx) error {
	app := WC.App{
		CustomerKey:    U.Env("WC_CONSUMER_KEY"),
		CustomerSecret: U.Env("WC_CONSUMER_SECRET"),
	}

	wc := WC.NewClient(app, U.Env("WC_SHOP_NAME"))

	// Retrieve the order details
	order, err := wc.Order.Get(int64(c.QueryInt("id")), nil)
	if err != nil {
		// log.Fatal(err)
		fmt.Printf("error: %v", err)
	}

	// Print the names of the products
	fmt.Println("Products purchased:")
	var productNames []string
	for _, item := range order.LineItems {
		fmt.Println(item.Name)
		productNames = append(productNames, item.Name)
	}
	return c.JSON(fiber.Map{"productNames": productNames})
}
