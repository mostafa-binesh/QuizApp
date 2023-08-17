package controllers

import (
	D "docker/database"
	M "docker/models"
	U "docker/utils"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AllStudyPlans(c *fiber.Ctx) error {
	studyPlans := []M.StudyPlan{}
	if err := D.DB().Find(&studyPlans).Error; err != nil {
		return U.DBError(c, err)
	}
	return c.JSON(fiber.Map{"data": studyPlans})
}
func CreateStudyPlanner(c *fiber.Ctx) error {
	payload := new(M.CreateNewStudyPlanInput)
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
	// Calculate the number of days between start and end dates
	numDays := int(payload.EndDate.Sub(payload.StartDate).Hours()/24) + 1

	// Get the starting day of the week (0 = Sunday, 1 = Monday, ..., 6 = Saturday)
	startingDay := int(payload.StartDate.Weekday())

	// Create study plans for each day within the date range
	for i := 0; i < numDays; i++ {
		date := payload.StartDate.Add(time.Duration(i) * 24 * time.Hour)
		workingHoursIndex := (startingDay + i) % len(payload.WorkingHours)
		hours := payload.WorkingHours[workingHoursIndex]
		plan := M.StudyPlan{
			Date:       date,
			Hours:      hours,
			IsFinished: false,
			UserID:     user.ID,
		}
		if err := D.DB().Create(&plan).Error; err != nil {
			return U.DBError(c, err)
		}
	}

	return c.JSON(fiber.Map{"msg": "Study plans created successfully"})
}
func VerifyDate(c *fiber.Ctx) error {
	payload := new(M.VerifyStudyPlanDateInput)
	// Parse the payload
	if err := c.BodyParser(payload); err != nil {
		U.ResErr(c, err.Error())
	}
	// Validsate the payload
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	// Get authenticated user
	user := c.Locals("user").(M.User)
	date := payload.Date
	// Update study plan
	result := D.DB().
		Model(&M.StudyPlan{}).
		Where("user_id = ? AND date = ?", user.ID, date).
		Update("is_finished", true)
	// Handle errors
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	if result.RowsAffected == 0 {
		return U.ResErr(c, fmt.Sprintf("Plan with date %s not found", date), 404)
	}
	return U.ResMessage(c, fmt.Sprintf("Plan with date %s has been verified", date))
}
