package controllers

import (
	D "docker/database"
	S "docker/services"

	// "math/rand"
	// "time"

	// F "docker/database/filters"
	M "docker/models"
	U "docker/utils"

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
	payload := new(M.AddCourseUsingOrderID)
	// parsing the payload
	if err := c.BodyParser(payload); err != nil {
		U.ResErr(c, err.Error())
	}
	// validation the payload
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	// get authenticated user
	user := c.Locals("user").(M.User)
	// get wc courses using payload.OrderID
	childCourses, purchasedCourseIDPayDateMap, err := S.ImportUserCoursesUsingOrderID(payload.OrderID)
	if err != nil {
		return U.DBError(c, err)
	}
	// convert childCourses to course_user model
	courseUser := S.AddCourseUserUsingCourses(childCourses, purchasedCourseIDPayDateMap, user.ID)
	// todo: we need to check if the course_user exist, then update the expiration date
	// we SAVE records to the database because some user may bought other courses already
	if err := D.DB().Save(&courseUser).Error; err != nil {
		return U.DBError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"msg": "Courses have been added"})
}
