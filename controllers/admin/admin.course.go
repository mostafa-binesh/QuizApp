package admin

import (
	D "docker/database"
	M "docker/models"
	S "docker/services"
	U "docker/utils"
	"github.com/gofiber/fiber/v2"
)
// todo: first need to insert all courses into the database and then again check for the parent's ids
func AddCoursesFromWooCommerce(c *fiber.Ctx) error {
	// get all woocommerce products from its service
	wcCourses, err := S.GetAllProducts()
	if err != nil {
		return U.ResErr(c, err.Error())
	}
	// convert the wooc. retreived products into course model
	convertedCourses, err := S.ConvertWCCourseToCourseModel(&wcCourses)
	if err != nil {
		return U.ResErr(c, err.Error())
	}
	// create a map to store parent courses id and its courses
	DBCourseWoocommerceIDMap := make(map[uint]*M.Course)
	for _, course := range *convertedCourses {
		// for some reason, some of courses were empty
		if course.Title == "" {
			continue
		}
		var existingCourse M.Course
		var parentCourse M.Course
		// Try to find the course with the given woocommerce_id
		result := D.DB().Where("woocommerce_id = ?", course.WoocommerceID).First(&existingCourse)
		// also need to find the parent course with the given course.parentid
		if course.ParentID != nil {
			if DBCourseWoocommerceIDMap[*course.ParentID] == nil {
				// parent course map doesn't exist in the map
				// try to find it in the database
				parentResult := D.DB().First(&parentCourse, *course.ParentID)
				if parentResult.RowsAffected == 1 {
					DBCourseWoocommerceIDMap[*course.ParentID] = &parentCourse
				}
			} else {
				parentCourse = *DBCourseWoocommerceIDMap[*course.ParentID]
			}
		}

		// if course with desired woocommerce_id not found, create it
		if result.RowsAffected == 0 {
			if parentCourse.Title != "" {
				course.ParentID = &parentCourse.ID
			}
			D.DB().Create(&course)
		} else {
			// if found, just update it
			existingCourse.Title = course.Title
			existingCourse.Duration = course.Duration
			if parentCourse.Title != "" {
				existingCourse.ParentID = &parentCourse.ID
			}
			existingCourse.ValidityDaysPeriod = course.ValidityDaysPeriod

			D.DB().Save(&existingCourse)
		}
	}
	return c.JSON(fiber.Map{"data": convertedCourses})
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
