package controllers

import (
	D "docker/database"
	S "docker/database/seeders"
	M "docker/models"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ! add any migration that you wanna add to the database
func AutoMigrate(c *fiber.Ctx) error {
	// ! drop all tables if 'dropAllTables' field is 1 in the query
	// return c.JSON(fiber.Map{
	// 	"message":c.Query("dropAllTables"),
	// })
	fmt.Println("dropAllTables")
	if c.Query("dropAllTables") == "1" {
		fmt.Println("dropping all tables")
		D.DB().Migrator().DropTable(
			&M.Course{},
			&M.User{},
			&M.Law{},
			&M.Comment{},
			&M.Keyword{},
			&M.File{},
			&M.Question{},
			&M.Option{},
			&M.UserAnswer{},
			&M.Subject{},
		)
	}
	fmt.Println("Tables migration done...")
	// ! migrate tables
	err := D.DB().AutoMigrate(
		&M.User{},
		&M.Subject{},
		&M.Law{},
		&M.Comment{},
		&M.Keyword{},
		&M.File{},
		&M.Course{},
		&M.Question{},
		&M.Option{},
		&M.UserAnswer{},
	)
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}
	// ! set seederRepeatCount, default 1
	var seederRepeatCount int64
	seederRepeatCount = 1
	seedCountQuery := c.Query("seederRepeatCount")
	if seedCountQuery != "" {
		var err error
		seederRepeatCount, err = strconv.ParseInt(seedCountQuery, 10, 64)
		if err != nil {
			panic("seedCount query param. cannot be parsed")
		}
	}
	// ! seeders
	fmt.Printf("seeder gonna run for %d loop", seederRepeatCount)
	for i := 0; i < int(seederRepeatCount); i++ {
		S.InitSeeder()
	}
	return c.SendString("migrate completed")
}
