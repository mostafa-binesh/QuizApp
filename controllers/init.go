package controllers

import (
	"fmt"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/fa"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	// en_translations "github.com/go-playground/validator/v10/translations/en"
	fa_translations "github.com/go-playground/validator/v10/translations/fa"

	// "github.com/gofiber/fiber/v2"
	D "docker/database"

	// "reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
	U "github.com/gofiber/fiber/v2/utils"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
	// en       *locales.Translator
)

func Initilize() error {
	// initilize validator translator
	en := en.New()
	fa := fa.New()
	uni = ut.New(en, fa)
	// trans, _ = uni.GetTranslator("en")
	trans, _ = uni.GetTranslator("fa")
	trans.Add("Description", "توضیحات", true)
	// trans.Add("توضیحات", "Description", false)
	// trans.Add("توضیحات", "Description", true)
	trans.Add("WebsiteURL", "user naee", true)
	trans.Add("Username", "user naee", true)
	validate = validator.New()
	fa_translations.RegisterDefaultTranslations(validate, trans)
	// ! add gorm exists and their translation, declaration is below
	validate.RegisterValidation("gexist", GormExists) // custom gorm validation
	validate.RegisterTranslation("gexist", trans, func(ut ut.Translator) error {
		return ut.Add("gexist", "{0} doesn't exist", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("gexist", fe.Field())

		return t
	})
	// ! add gorm unique and their translation, declaration is below
	validate.RegisterValidation("gunique", GormUnique) // custom gorm validation
	validate.RegisterTranslation("gunique", trans, func(ut ut.Translator) error {
		return ut.Add("gunique", "{0} already exists!", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("gunique", fe.Field())

		return t
	})
	return nil
}

// ! gorm exist function
func GormExists(fl validator.FieldLevel) bool {
	params := strings.Split(fl.Param(), ".")
	// 2 params like exists,tableName and the columnName should be retreieved from field's name
	// 3 params like exists,tableName,ColumnName and everything should be put in the query in the SAFE MODE !
	// exists and unique are two options that can be called in the code
	var query string
	switch len(params) {
	case 1:
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", params[0], U.ToLower(fl.FieldName()))
	case 2:
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", params[0], params[1])
	}
	fmt.Printf("query is: %s", query)
	fmt.Printf("| field string value is : %s", fl.Field().String())
	rowsCount := D.RowsCount(query, fl.Field().String())
	fmt.Printf("rowsCount: %v", rowsCount)
	return rowsCount > 0
}

// ! gorm unique function
func GormUnique(fl validator.FieldLevel) bool {
	params := strings.Split(fl.Param(), ".")
	// 2 params like exists,tableName and the columnName should be retreieved from field's name
	// 3 params like exists,tableName,ColumnName and everything should be put in the query in the SAFE MODE !
	// exists and unique are two options that can be called in the code
	var query string
	switch len(params) {
	case 1:
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", params[0], U.ToLower(fl.FieldName()))
	case 2:
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", params[0], params[1])
	}
	fmt.Printf("query is: %s", query)
	fmt.Printf("| field string value is : %s", fl.Field().String())
	rowsCount := D.RowsCount(query, fl.Field().String())
	fmt.Printf("rowsCount: %v", rowsCount)
	return rowsCount == 0
	fmt.Println("until last false")
	return false
}

// ! this function is being used for returned good-reading json errors
func ValidatorErrors(err validator.ValidationErrorsTranslations) map[string]string {
	// Define variable for error fields.
	errFields := map[string]string{}
	for k, v := range err {
		// fmt.Printf("key[%s] value[%s]\n", k, v)
		errFields[strings.Split(k, ".")[1]] = v
	}
	return errFields
}

// ! return error if any error exist
// ! use ValidateErrors function
func ValidationHandle(c *fiber.Ctx, err error) error {
	// translate all error at once
	errs := err.(validator.ValidationErrors)
	// returns a map with key = namespace & value = translated error
	// NOTICE: 2 errors are returned and you'll see something surprising
	// translations are i18n aware!!!!
	// eg. '10 characters' vs '1 character'
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"errors": ValidatorErrors(errs.Translate(trans)),
	})
}

// ! handling BodyParser, returns two value, response and error
// ! how to work with it: res, err := BodyParserHandle(c,user) , if err != nil return res
func BodyParserHandle(c *fiber.Ctx, payload interface{}) (error, error) {
	if err := c.BodyParser(payload); err != nil {
		fmt.Println("bodyParserHandle function error")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()}), err
	}
	return nil, nil
}

// ! return a custom error
func ReturnError(c *fiber.Ctx, err string, statusCode ...int) error {
	x := fiber.StatusBadRequest
	if len(statusCode) > 0 {
		x = statusCode[0]
	}
	return c.Status(x).JSON(fiber.Map{"message": err})
}

// ! some of the init.go functions don't belong in here, make a better structure for them
