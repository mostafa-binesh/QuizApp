package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	// "github.com/go-playground/locales/en"
	D "docker/database"
	// "gorm.io/gorm"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/fa"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	fa_translations "github.com/go-playground/validator/v10/translations/fa"
	"github.com/iancoleman/strcase"
	// en_translations "github.com/go-playground/validator/v10/translations/en"
)

var faTranslation = map[string]string{
	"Name":               "نام",
	"FirstName":          "نام",
	"FullName":           "نام",
	"Email":              "ایمیل",
	"LastName":           "نام خانوادگی",
	"File":               "فایل",
	"Username":           "نام کاربری",
	"Password":           "رمز عبور",
	"Age":                "سن",
	"Type":               "نوع",
	"Title":              "عنوان",
	"SessionNumber":      "شماره جلسه",
	"SessionDate":        "تاریخ جلسه",
	"NotificationNumber": "شماره ابلاغ",
	"NotificationDate":   "تاریخ ابلاغ",
	"Body":               "بدنه",
	"Tags":               "تگ ها",
	"Image":              "عکس",
	"PersonalCode":       "کد پرسنلی",
	"NationalCode":       "کد ملی",
	"PhoneNumber":        "شماره همراه",
	"Category":           "دسته بندی",
	"OrderID":            "شماره سفارش",
}
var IgnoreID uint64

// gets fields as interface and ignoreID if you want to ignore a certain id
func Validate(fields interface{}, ignoreID ...string) map[string]string {
	if len(ignoreID) > 0 {
		ignoreIDUint64, err := strconv.ParseUint(ignoreID[0], 10, 64)
		if err != nil {
			panic("validate function, cannot parse ignoreID")
		}
		IgnoreID = ignoreIDUint64
	}
	en := en.New()
	fa := fa.New()
	uni := ut.New(en, fa)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ := uni.GetTranslator("fa")
	validate := validator.New()
	fa_translations.RegisterDefaultTranslations(validate, trans)
	// ! custom names registration
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return faTranslation[field.Name]
	})
	// ! custom translations
	// ? required
	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} الزامی است", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})
	// ? gorm unique
	if err := validate.RegisterValidation("dunique", uniqueValidator); err != nil {
		panic(err)
	}
	validate.RegisterTranslation("dunique", trans, func(ut ut.Translator) error {
		return ut.Add("dunique", "{0} ثبت شده است", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("dunique", fe.Field())
		return t
	})
	// ? gorm exists
	if err := validate.RegisterValidation("dexists", existValidator); err != nil {
		panic(err)
	}
	validate.RegisterTranslation("dexists", trans, func(ut ut.Translator) error {
		return ut.Add("dexists", "{0} ثبت نشده است", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("dexists", fe.Field())
		return t
	})
	// ! possible issues: if fields have another struct in it, getJSONTag
	// ! > won't work properly
	err := validate.Struct(fields)
	if err != nil {
		responseError := make(map[string]string)
		errs := err.(validator.ValidationErrors)
		var jsonTag string
		for _, e := range errs {
			jsonTag = GetJSONTag(fields, e.StructField())
			if jsonTag == "" {
				jsonTag = ToLowerCamel(e.StructField())
			}
			responseError[jsonTag] = e.Translate(trans) // works fine
		}
		return responseError
	}
	return nil
}

// func uniqueValidator(fl validator.FieldLevel, ignore ...uint) func(fl validator.FieldLevel) bool {
func uniqueValidator(fl validator.FieldLevel) bool {
	// ! validator should be like "dunique:<tableName>.<Optional:columnName>" like "dunique:users.email", "dunique:users" (if field is phoneNumber, it will be converted to "dunique:users.phone_number")
	// ! if you wanna ignore a id (eg. in edit) you need to pass the id to second argument
	// ! -- of validation function
	fmt.Printf("tag: %s\n", fl.GetTag())
	fmt.Printf("field: %s\n", fl.Field())
	fmt.Printf("fieldName: %s\n", fl.FieldName())
	fmt.Printf("parent: %s\n", fl.Parent())
	fmt.Printf("struct field name: %s\n", fl.StructFieldName())
	fmt.Printf("param: %s\n", fl.Param())
	fmt.Printf("field: %s\n", fl.Field())
	params := strings.Split(fl.Param(), ".")
	fmt.Printf("params: %v\n", fl.Param())
	fmt.Printf("parsed params: %v\n", params)
	// 2 params like exists,tableName and the columnName should be retreieved from field's name
	// 3 params like exists,tableName,ColumnName and everything should be put in the query in the SAFE MODE !
	// exists and unique are two options that can be called in the code
	var query string
	switch len(params) {
	case 1:
		query = fmt.Sprintf("SELECT id FROM %s WHERE %s = ?", params[0], strcase.ToSnake(fl.StructFieldName()))
	case 2:
		query = fmt.Sprintf("SELECT id FROM %s WHERE %s = ?", params[0], params[1])
	}
	fmt.Printf("query is: %s", query)
	fmt.Printf("| field string value is : %s", fl.Field().String())
	rowsCount := D.RowsCount(query, fl.Field().String(), IgnoreID)
	fmt.Printf("rowsCount: %d\n", rowsCount)
	return rowsCount == 0
}

// }
func existValidator(fl validator.FieldLevel) bool {
	return !uniqueValidator(fl)
}
