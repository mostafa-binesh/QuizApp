package utils

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// returns status 400 of { "error": DBError}
func DBError(c *fiber.Ctx, err error) error {
	var errorText string
	if errors.Is(err, gorm.ErrRecordNotFound) {
		errorText = "Data not found"
	} else if errors.Is(err, gorm.ErrInvalidData) {
		errorText = "Invalid Data"
	} else if errors.Is(err, gorm.ErrDuplicatedKey) {
		errorText = "Duplicate Key Error"
	} else {
		errorText = "Unpredicted Database Error"
	}
	fmt.Printf("Database error: %v\n", err)
	if Env("APP_DEBUG") == "true" {
		return c.Status(400).JSON(fiber.Map{
			"error": errorText,
			"debug": err,
		})
	}
	return ResErr(c, errorText)
}
