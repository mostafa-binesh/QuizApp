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
	// ! order doesn't matter in drop and creating the tables
	// ! > but it does matter in seeding
	// ! drop all tables if 'dropAllTables' field is 1 in the query
	fmt.Println("dropAllTables")
	if c.Query("dropAllTables") == "1" {
		fmt.Println("dropping all tables")
		// delete many to many tables first
		D.DB().Exec("DELETE FROM course_user")
		D.DB().Exec("DELETE FROM user_answers")
		D.DB().Migrator().DropTable(
			&M.User{},
			&M.Course{},
			&M.Question{},
			&M.Option{},
			&M.UserAnswer{},
			&M.System{},
			&M.Subject{},
			&M.Quiz{},
			&M.Image{},
			&M.Tab{},
			&M.Dropdown{},
			&M.CourseUser{},
			&M.StudyPlan{},
		)
	}
	if c.QueryInt("justDrop") == 1 {
		return c.SendString("operation dropped by you")
	}
	fmt.Println("Tables migration done...")
	// ! migrate tables
	err := D.DB().AutoMigrate(
		&M.User{},
		&M.Course{},
		&M.Quiz{},
		&M.UserAnswer{},
		&M.Subject{},
		&M.System{},
		&M.Tab{},
		&M.Question{},
		&M.Option{},
		&M.Image{},
		&M.Dropdown{},
		&M.CourseUser{},
		&M.StudyPlan{},
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
