package utils

import (
	"reflect"
)

func StructToMap(in interface{}) map[string]interface{} {
	v := reflect.ValueOf(in)
	vType := v.Type()

	result := make(map[string]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		name := vType.Field(i).Tag.Get("json")
		result[name] = v.Field(i).Interface()
	}

	return result
}
