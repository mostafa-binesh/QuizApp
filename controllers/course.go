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
	user := M.AuthedUser(c)
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
	user := M.AuthedUser(c)
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
	user := M.AuthedUser(c)
	// get wc courses using payload.OrderID
	childCourses, purchasedCourseIDPayDateMap, err := S.ImportUserCoursesUsingOrderID(payload.OrderID)
	if err != nil {
		return U.DBError(c, err)
	}
	// get user's bought courses with orderID and seperate them to already-inserted and newly courses
	courseUsersToUpdate, newCourseUsers, err := S.ExtractCourseToInsertAndToUpdate(childCourses, purchasedCourseIDPayDateMap, user.ID)
	if err != nil {
		return U.DBError(c, err)
	}
	// Batch update existing courseUser records
	if len(courseUsersToUpdate) > 0 {
		if err := D.DB().Save(courseUsersToUpdate).Error; err != nil {
			return U.DBError(c, err)
		}
	}
	// Batch insert new courseUser records
	if len(newCourseUsers) > 0 {
		if err := D.DB().Create(&newCourseUsers).Error; err != nil {
			return U.DBError(c, err)
		}
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"msg": "Courses have been added"})
}
