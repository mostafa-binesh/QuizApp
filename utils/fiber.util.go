package utils

import "github.com/gofiber/fiber/v2"

var fiberContext *fiber.Ctx

func SetFiberContext(c *fiber.Ctx) {
	fiberContext = c
}
func FiberCtx() *fiber.Ctx {
	return fiberContext
}
