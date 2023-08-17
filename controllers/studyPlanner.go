package controllers

import (
	D "docker/database"
	M "docker/models"
	U "docker/utils"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

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
	// Find the desired plan which belongs to this user
	studyPlan := M.StudyPlan{}
	result := D.DB().
		Where("user_id = ? AND date = ?", user.ID, payload.Date).
		First(&studyPlan)
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	// Finish the study plan
	studyPlan.Finish()
	if err := D.DB().Save(&studyPlan).Error; err != nil {
		return U.DBError(c, err)
	}
	// Show error if no plan found with desired date
	if result.RowsAffected == 0 {
		return U.ResErr(c, fmt.Sprintf("Plan with date %s not found", payload.Date))
	}
	return U.ResMessage(c, fmt.Sprintf("Plan with date %s has been verified", payload.Date))
}
