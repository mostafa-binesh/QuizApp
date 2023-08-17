package controllers

import (
	D "docker/database"
	M "docker/models"
	U "docker/utils"
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
	// Create study plans for each day within the date range
	for i := 0; i < numDays; i++ {
		date := payload.StartDate.Add(time.Duration(i) * 24 * time.Hour)
		for _, hours := range payload.WorkingHours {
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
	}

	return c.JSON(fiber.Map{"message": "Study plans created successfully"})
}
