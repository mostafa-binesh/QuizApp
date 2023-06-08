package utils

import (
	"reflect"
)

func GetJSONTag(s interface{}, field string) string {
	// get the reflect.Value of the struct
	v := reflect.ValueOf(s)

	// check that the given interface is a struct
	if v.Kind() != reflect.Struct {
		return ""
	}

	// iterate over the fields of the struct
	for i := 0; i < v.NumField(); i++ {
		// get the reflect.Type of the field
		fieldType := v.Type().Field(i)
		// check if the field name matches the desired field
		if fieldType.Name == field {
			jsonTag := fieldType.Tag.Get("json")
			// return the tag value
			return jsonTag
		}
	}

	// return an empty string if the field is not found
	return ""
}
