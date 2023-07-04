package admin
// this controller hasn't been used in this project 
// it's here just because has some code templates 
import (
	D "docker/database"
	F "docker/database/filters"
	U "docker/utils"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	// "strconv"

	// F "docker/database/filters"
	M "docker/models"
)

func IndexLaw(c *fiber.Ctx) error {
	laws := []M.Law{}
	pagination := U.ParsedPagination(c)
	D.DB().Scopes(U.Paginate(laws, pagination)).Find(&laws)
	responseLaws := []M.LawMinimal_min{}
	for i := 0; i < len(laws); i++ {
		responseLaws = append(responseLaws, M.LawMinimal_min{
			ID:    laws[i].ID,
			Title: laws[i].Title,
			Image: laws[i].Image,
		})
	}
	return U.ResWithPagination(c, responseLaws, *pagination)
}
func UpdateLaw(c *fiber.Ctx) error {
	law := M.Law{}
	payload := new(M.CreateLawInput)
	// parsing the payload
	if err := c.BodyParser(payload); err != nil {
		U.ResErr(c, err.Error())
	}
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	result1 := D.DB().Where("id = ?", c.Params("id")).Find(&law)
	if result1.Error != nil {
		return U.DBError(c, result1.Error)
	}
	law.Body = payload.Body
	law.Image = payload.Image
	law.NotificationDate = payload.NotificationDate
	law.NotificationNumber = payload.NotificationNumber
	law.SessionDate = payload.SessionDate
	law.SessionNumber = payload.SessionNumber
	law.Title = payload.Title
	law.Type = payload.Type
	result := D.DB().Save(&law)
	if result.Error != nil {
		return U.ResErr(c, "مشکلی در به روز رسانی به وجود آمده")
	}
	// ! store ExplanatoryPlan is exists
	file, _ := c.FormFile("explanatoryPlan")
	// if formError != nil {
	// 	return U.ResErr(c, formError.Error())
	// }
	if file != nil {
		fmt.Println("till here")
		// check if file with this name already exists
		if U.FileExistenceCheck(file.Filename, "./public/uploads") {
			return U.ResErr(c, "file already exists")
		}
		// ! file extension check
		// if !(U.HasImageSuffixCheck(file.Filename) || U.HasSuffixCheck(file.Filename, []string{"pdf"})) {
		// 	return c.SendString("file should be image or pdf! please fix it")
		// }
		// Save file to disk
		fileName := U.AddUUIDToString(file.Filename)
		c.SaveFile(file, fmt.Sprintf("./public/uploads/%s", fileName))
		D.DB().Create(&M.File{
			Type:  M.FileTypes["plan"],
			Name:  fileName,
			LawID: law.ID,
		})
	}
	// ! certificate
	file, _ = c.FormFile("certificate")
	if file != nil {
		// check if file with this name already exists
		if U.FileExistenceCheck(file.Filename, "./public/uploads") {
			return U.ResErr(c, "file already exists")
		}
		// ! file extension check
		// if !(U.HasImageSuffixCheck(file.Filename) || U.HasSuffixCheck(file.Filename, []string{"pdf"})) {
		// 	return c.SendString("file should be image or pdf! please fix it")
		// }
		// Save file to disk
		fileName := U.AddUUIDToString(file.Filename)
		c.SaveFile(file, fmt.Sprintf("./public/uploads/%s", fileName))
		D.DB().Create(&M.File{
			Type:  M.FileTypes["certificate"],
			Name:  fileName,
			LawID: law.ID,
		})
	}
	// ! attachments
	// attachments, _ := c.FormFile("explanatoryPlan")
	form, _ := c.MultipartForm()
	attachments := form.File["attachments"]
	for _, file := range attachments {
		// check if file with this name already exists
		if U.FileExistenceCheck(file.Filename, "./public/uploads") {
			return U.ResErr(c, "file already exists")
		}
		// ! file extension check
		// if !(U.HasImageSuffixCheck(file.Filename) || U.HasSuffixCheck(file.Filename, []string{"pdf"})) {
		// 	return c.SendString("file should be image or pdf! please fix it")
		// }
		// Save file to disk
		// err = c.SaveFile(file, fmt.Sprintf("./public/uploads/%s", file.Filename))
		fileName := U.AddUUIDToString(file.Filename)
		c.SaveFile(file, fmt.Sprintf("./public/uploads/%s", fileName))
		D.DB().Create(&M.File{
			Type:  M.FileTypes["attachment"],
			Name:  fileName,
			LawID: law.ID,
		})
		// if err != nil {
		// 	return U.ResErr(c, "cannot save")
		// }
	}
	return c.JSON(fiber.Map{
		"msg": "به روز رسانی با موفقیت انجام شد",
	})
}
func DeleteLaw(c *fiber.Ctx) error {
	result := D.DB().Delete(&M.Law{}, c.Params("id"))
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	if result.RowsAffected == 0 {
		return U.ResErr(c, "مصوبه یافت نشد")
	}
	return c.JSON(fiber.Map{
		"msg": "حذف کردن با موفقیت انجام شد",
	})
}
func CreateLaw(c *fiber.Ctx) error {
	payload := new(M.CreateLawInput)
	// parse payload
	if err := c.BodyParser(payload); err != nil {
		U.ResErr(c, err.Error())
	}
	// validate payload
	if errs := U.Validate(payload); errs != nil {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
	}
	law := M.Law{
		Type:               payload.Type,
		Title:              payload.Title,
		SessionNumber:      payload.SessionNumber,
		SessionDate:        payload.SessionDate,
		NotificationDate:   payload.NotificationDate,
		NotificationNumber: payload.NotificationNumber,
		Body:               payload.Body,
		Image:              payload.Image,
	}
	// store law in the db
	result := D.DB().Create(&law)
	if result.Error != nil {
		return U.ResErr(c, result.Error.Error())
	}
	// parse tags and store them
	var tags = strings.Split(payload.Tags, ",")
	for i := 0; i < len(tags); i++ {
		result2 := D.DB().Create(&M.Keyword{
			Keyword: tags[i],
			LawID:   law.ID,
		})
		if result2.Error != nil {
			D.DB().Delete(&M.Law{}, law.ID)
			return U.ResErr(c, result.Error.Error())
		}
	}
	// ! store ExplanatoryPlan is exists
	file, _ := c.FormFile("explanatoryPlan")
	// if formError != nil {
	// 	return U.ResErr(c, formError.Error())
	// }
	if file != nil {
		fmt.Println("till here")
		// check if file with this name already exists
		if U.FileExistenceCheck(file.Filename, "./public/uploads") {
			return U.ResErr(c, "file already exists")
		}
		// ! file extension check
		// if !(U.HasImageSuffixCheck(file.Filename) || U.HasSuffixCheck(file.Filename, []string{"pdf"})) {
		// 	return c.SendString("file should be image or pdf! please fix it")
		// }
		// Save file to disk
		fileName := U.AddUUIDToString(file.Filename)
		c.SaveFile(file, fmt.Sprintf("./public/uploads/%s", fileName))
		D.DB().Create(&M.File{
			Type:  M.FileTypes["plan"],
			Name:  fileName,
			LawID: law.ID,
		})
	}
	// ! certificate
	file, _ = c.FormFile("certificate")
	if file != nil {
		// check if file with this name already exists
		if U.FileExistenceCheck(file.Filename, "./public/uploads") {
			return U.ResErr(c, "file already exists")
		}
		// ! file extension check
		// if !(U.HasImageSuffixCheck(file.Filename) || U.HasSuffixCheck(file.Filename, []string{"pdf"})) {
		// 	return c.SendString("file should be image or pdf! please fix it")
		// }
		// Save file to disk
		fileName := U.AddUUIDToString(file.Filename)
		c.SaveFile(file, fmt.Sprintf("./public/uploads/%s", fileName))
		D.DB().Create(&M.File{
			Type:  M.FileTypes["certificate"],
			Name:  fileName,
			LawID: law.ID,
		})
	}
	// ! attachments
	// attachments, _ := c.FormFile("explanatoryPlan")
	form, _ := c.MultipartForm()
	attachments := form.File["attachments"]
	for _, file := range attachments {
		// check if file with this name already exists
		if U.FileExistenceCheck(file.Filename, "./public/uploads") {
			return U.ResErr(c, "file already exists")
		}
		// ! file extension check
		// if !(U.HasImageSuffixCheck(file.Filename) || U.HasSuffixCheck(file.Filename, []string{"pdf"})) {
		// 	return c.SendString("file should be image or pdf! please fix it")
		// }
		// Save file to disk
		// err = c.SaveFile(file, fmt.Sprintf("./public/uploads/%s", file.Filename))
		fileName := U.AddUUIDToString(file.Filename)
		c.SaveFile(file, fmt.Sprintf("./public/uploads/%s", fileName))
		D.DB().Create(&M.File{
			Type:  M.FileTypes["attachment"],
			Name:  fileName,
			LawID: law.ID,
		})
		// if err != nil {
		// 	return U.ResErr(c, "cannot save")
		// }
	}
	// return response
	return c.Status(200).JSON(fiber.Map{
		"msg": "مصوبه با موفقیت اضافه شد",
	})
}
func LawSearch(c *fiber.Ctx) error {
	laws := []M.Law{}
	pagination := U.ParsedPagination(c)
	D.DB().Scopes(
		F.FilterByType(c,
			F.FilterType{QueryName: "title", Operator: "LIKE"},
			F.FilterType{QueryName: "session_number", Operator: "LIKE"},
			F.FilterType{QueryName: "notification_number", Operator: "LIKE"},
			F.FilterType{QueryName: "body", Operator: "LIKE"},
			F.FilterType{QueryName: "recommender", Operator: "LIKE"},
			F.FilterType{QueryName: "type", Operator: "="},
			F.FilterType{QueryName: "notification_startDate", ColumnName: "notification_date", Operator: ">="},
			F.FilterType{QueryName: "notification_endDate", ColumnName: "notification_date", Operator: "<="},
			F.FilterType{QueryName: "session_startDate", ColumnName: "session_date", Operator: ">="},
			F.FilterType{QueryName: "session_endDate", ColumnName: "session_date", Operator: "<="}),
		U.Paginate(laws, pagination)).
		Find(&laws)
	pass_data := []M.LawMinimal_min{}
	for i := 0; i < len(laws); i++ {
		pass_data = append(pass_data, M.LawMinimal_min{
			ID:    laws[i].ID,
			Title: laws[i].Title,
			Image: laws[i].Image,
		})
	}
	return c.JSON(fiber.Map{
		"meta": pagination,
		"data": pass_data,
	})
}
func LawByID(c *fiber.Ctx) error {
	law := &M.Law{}
	if err := D.DB().Preload("Comments.User").Preload("Files").First(law, c.Params("id")).Error; err != nil {
		return U.DBError(c, err)
	}
	LawByID := M.LawToLawByID(law)
	return c.JSON(fiber.Map{
		"data": LawByID,
	})
}
func DeleteFile(c *fiber.Ctx) error {
	result := D.DB().Delete(&M.File{}, c.Params("fileID"))
	if result.Error != nil {
		return U.DBError(c, result.Error)
	}
	if result.RowsAffected == 0 {
		return U.ResErr(c, "فایل یافت نشد")
	}
	return c.JSON(fiber.Map{
		"msg": "فایل حذف شد",
	})
}
