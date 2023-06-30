package admin

import (
	D "docker/database"
	M "docker/models"
	S "docker/services"
	U "docker/utils"
	"github.com/gofiber/fiber/v2"
)

func AddCoursesFromWooCommerce(c *fiber.Ctx) error {
	// get all woocommerce products from its service
	courses, err := S.GetAllProducts()
	if err != nil {
		return U.ResErr(c, err.Error())
	}
	// create every course in the database
	for _, course := range courses {
		D.DB().Create(&M.Course{
			WoocommerceID: uint(course.ID),
			Title:         course.Name,
		})
	}
	return c.JSON(fiber.Map{"data": courses})
}
func AllCourses(c *fiber.Ctx) error {
	// get all courses
	courses := []M.Course{}
	result := D.DB().Preload("Subjects.Systems").Find(&courses)
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	// convert course to courseWithTitleOnly
	var CoursesWithTitleOnly []M.CourseWithTitleOnly
	for _, course := range courses {
		CoursesWithTitleOnly = append(CoursesWithTitleOnly, M.CourseWithTitleOnly{
			ID:    course.ID,
			Title: course.Title,
		})
	}
	// return courses
	return c.JSON(fiber.Map{"data": CoursesWithTitleOnly})
}
func CourseByID(c *fiber.Ctx) error {
	course := &M.Course{}
	result := D.DB().Find(course, c.Params("id"))
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	return c.JSON(fiber.Map{"data": course})
}
func CreateCourse(c *fiber.Ctx) error {
	payload := new(M.CourseInput)
	// parsing the payload
	if err := c.BodyParser(payload); err != nil {
		U.ResErr(c, err.Error())
	}
	// validate the payload
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	// craete the course into the database
	result := D.DB().Create(&M.Course{
		Title:         payload.Title,
		WoocommerceID: payload.WoocommerceID,
	})
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	// return success message
	return U.ResMessage(c, "دوره اضافه شد")
}
func UpdateCourse(c *fiber.Ctx) error {
	course := M.Course{}
	payload := new(M.CourseInput)
	// parsing the payload
	if err := c.BodyParser(payload); err != nil {
		U.ResErr(c, err.Error())
	}
	// validate the payload
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	result := D.DB().Where("id = ?", c.Params("id")).Find(&course)
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	course.Title = payload.Title
	course.WoocommerceID = payload.WoocommerceID
	result = D.DB().Save(&course)
	if result.Error != nil {
		return U.ResErr(c, "مشکلی در به روز رسانی به وجود آمده")
	}
	return U.ResMessage(c, "دوره بروز شد")
}

// returns all courses.subjects.systems
func AllSubjects(c *fiber.Ctx) error {
	courses := []M.Course{}
	result := D.DB().Preload("Subjects.Systems").Find(&courses)
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	return c.JSON(fiber.Map{"data": courses})
}

// returns all subjects.systems of course with param of courseID
func CourseSubjects(c *fiber.Ctx) error {
	// get user's course with id of param with subject with system of the user
	course := M.Course{}
	result := D.DB().Preload("Subjects.Systems").First(&course, c.Params("courseID"))
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	if result.RowsAffected != 1 {
		return U.ResErr(c, "Course not found")
	}
	// change subject model to subjectWithSystems model
	subjectWithSystems := []M.SubjectWithSystems{}
	for j := 0; j < len(course.Subjects); j++ {
		subjectWithSystems = append(subjectWithSystems, M.SubjectWithSystems{
			ID:       course.Subjects[j].ID,
			Title:    course.Subjects[j].Title,
			Systems:  course.Subjects[j].Systems,
			CourseID: course.ID,
		})
	}
	return c.JSON(fiber.Map{"data": subjectWithSystems})
}
