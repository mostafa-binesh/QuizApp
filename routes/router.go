package routes

import (
	U "docker/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func RouterInit() {
	router := fiber.New(fiber.Config{
		ServerHeader: "ubirockteam@gmail.com",
		AppName:      "Medical Exam Quiz Application",
	})
	// ! add middleware
	// cors
	router.Use(cors.New(cors.Config{
		AllowOrigins:     U.Env("APP_ALLOW_ORIGINS"),
		AllowCredentials: true,
	}))
	// setup utility base url
	router.Use(func(c *fiber.Ctx) error {
		U.BaseURL = c.BaseURL()
		return c.Next()
	})
	// setup fiber context utility
	router.Use(func(c *fiber.Ctx) error {
		U.SetFiberContext(c)
		return c.Next()
	})
	// logger
	router.Use(logger.New())
	// recovery from panic
	router.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	// #######################
	router.Static("/public", "./public") // static files, local public folder to public url

	// ! api routes
	APIInit(router)
	// ! listen
	router.Listen(":" + U.Env("APP_PORT"))
	// if U.Env("environment") == "development" {
	// 	router.Listen("localhost:" + U.Env("APP_PORT"))
	// } else if U.Env("environment") == "production" {
	// 	router.Listen(":" + U.Env("APP_PORT"))
	// }
}
