package utils

import (
	// F "docker/database/filters"
	"github.com/gofiber/fiber/v2"
)

// return status 400 of { "error" : errMessage }
func ResErr(c *fiber.Ctx, err string, statusCode ...int) error {
	status := fiber.StatusBadRequest
	if len(statusCode) > 0 {
		status = statusCode[0]
	}
	return c.Status(status).JSON(fiber.Map{"error": err})
}
func ResDebug(c *fiber.Ctx, err error, errorOptionalText ...string) error {
	errorText := "Internal Error"
	if len(errorOptionalText) > 0 {
		errorText = errorOptionalText[0]
	}
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errorText, "debug": err})
}
func ResValidationErr(c *fiber.Ctx, err map[string]string) error {
	return c.Status(400).JSON(fiber.Map{
		"errors": err,
	})
}
func ResWithPagination(c *fiber.Ctx, data interface{}, pagination Pagination) error {
	return c.Status(200).JSON(fiber.Map{
		"meta": pagination,
		"data": data,
	})
}

// returns status 200 of { "msg" : sendMessage }
func ResMsg(c *fiber.Ctx, msg string, statusCode ...int) error {
	status := fiber.StatusOK
	if len(statusCode) > 0 {
		status = statusCode[0]
	}
	return c.Status(status).JSON(fiber.Map{
		"msg": msg,
	})
}
