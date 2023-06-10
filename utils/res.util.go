package utils

import (
	// F "docker/database/filters"
	"github.com/gofiber/fiber/v2"
)

// response error, easier way to return a json error
func ResErr(c *fiber.Ctx, err string) error {
	// return FiberCtx().Status(400).JSON(fiber.Map{
	return c.Status(400).JSON(fiber.Map{
		"error": err,
	})
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

func ResMessage(c *fiber.Ctx, msg string) error {
	return c.Status(200).JSON(fiber.Map{
		"msg": msg,
	})
}
