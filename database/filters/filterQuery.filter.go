package filters

import (
	// "docker/utils"
	"github.com/iancoleman/strcase"

	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// TODO : convert all foreach to for loops
type FilterType struct {
	QueryName  string
	ColumnName string
	Operator   string // eg. LIKE, BETWEEN and etc.
}

// key and values
func FilterByMap(queryParams map[string]string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for key, value := range queryParams {
			if value == "" {
				continue
			}
			db = db.Where(fmt.Sprintf("%s = ?", key), value)
		}
		return db
	}
}

// fiber context and parameters, eg. name, title
func FilterByParameters(c *fiber.Ctx, queryParams []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var queryValue string
		for _, value := range queryParams {
			queryValue = c.Query(value)
			if queryValue == "" {
				continue
			}
			fmt.Printf("%v", queryValue)
			db = db.Where(fmt.Sprintf("%s = ?", value), queryValue)
		}
		return db
	}
}

// interface should have column tag, if not, snakeCase of the name would be consider
// interface also can have a operator tag, if not, operator will be assigned as default (=)
func FilterByInterface(c *fiber.Ctx, u interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		v := reflect.ValueOf(u)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		var queryValue string
		var operator string
		var jsonTag string
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			if jsonTag = field.Tag.Get("column"); jsonTag == "" {
				jsonTag = strcase.ToLowerCamel(field.Name)
			}
			value := v.Field(i).Interface()
			queryValue = c.Query(jsonTag)
			// ! if there is no value for query parameter, continue the loop
			if queryValue == "" {
				continue
			}
			if operator = field.Tag.Get("operator"); operator == "" {
				operator = "="
			}
			db = db.Where(fmt.Sprintf("%s %s ?", jsonTag, operator), queryValue)
			println(jsonTag, ":", value)
		}
		return db
	}
}

// ! eg. filterType{ QueryName: "name"} >> there should be a name parameter in the query
// ! eg. filterType{ QueryName: "body", Operator: "LIKE"} >> there should be a body parameter in the query
// ! eg. filterType{ QueryName: "startDate",ColumnName: "release_date", Operator: ">="}
// ! columnName is optional, if not exist, queryName will be considered as columnName
// ! Operator is optional as well, default is =
func FilterByType(c *fiber.Ctx, filterTypes ...FilterType) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var value string
		var queryValue string
		for _, filter := range filterTypes {
			queryValue = c.Query(filter.QueryName)
			if queryValue == "" {
				continue
			}
			if filter.ColumnName == "" {
				filter.ColumnName = filter.QueryName
			}
			if filter.Operator == "LIKE" {
				value = fmt.Sprintf("%%%s%%", queryValue)
			} else {
				value = queryValue
			}
			db = db.Where(fmt.Sprintf("%s %s ?", filter.ColumnName, filter.Operator), value)
		}
		return db
	}
}
