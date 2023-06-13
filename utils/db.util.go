package utils

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DBError(c *fiber.Ctx, err error) error {
	var errorText string
	if errors.Is(err, gorm.ErrRecordNotFound) {
		errorText = "داده یافت نشد"
	} else if errors.Is(err, gorm.ErrInvalidData) {
		errorText = "داده نامعتبر است"
	} else if errors.Is(err, gorm.ErrDuplicatedKey) {
		errorText = "مقدار تکراری در پایگاه داده وجود دارد"
	} else {
		errorText = "خطای پیش بینی نشده ی پایگاه داده"
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
