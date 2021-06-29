package utils

import (
	"reflect"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
func StructToMap(in interface{}) map[string]interface{} {
	v := reflect.ValueOf(in)
	vType := v.Type()

	result := make(map[string]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		name := vType.Field(i).Name

		result[ToSnakeCase(name)] = v.Field(i).Interface()
	}

	return result
}
