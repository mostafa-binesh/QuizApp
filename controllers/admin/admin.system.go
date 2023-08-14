package admin

import (
	D "docker/database"
	M "docker/models"
	U "docker/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// returns all questions with systemID of parameter "systemID"
func SystemQuestions(c *fiber.Ctx) error {
	var questions []M.Question
	if err := D.DB().
		Where("system_id = ?", c.Params("systemID")).
		Preload("Options").
		Find(&questions).
		Error; err != nil {
		return U.DBError(c, err)
	}
	return c.JSON(fiber.Map{"data": questions})
}

// deletes all questions with systemID of parameter "systemID"
func DeleteSystemQuestions(c *fiber.Ctx) error {
	if err := D.DB().Where("system_id = ?", c.Params("systemID")).Delete(&M.Question{}).Error; err != nil {
		return U.DBError(c, err)
	}
	return c.JSON(fiber.Map{"msg": fmt.Sprintf("All Questions with SystemID of %s have been deleted", c.Params("systemID"))})
}
