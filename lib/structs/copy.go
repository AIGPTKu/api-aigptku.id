package structs

import "reflect"

func Copy(src any, dest any) {
	srcValue := reflect.ValueOf(src)
	destValue := reflect.ValueOf(dest).Elem() // Dereference the pointer

	for i := 0; i < destValue.NumField(); i++ {
		destFieldName := destValue.Type().Field(i).Name
		srcFieldValue := srcValue.FieldByName(destFieldName)
		if srcFieldValue.IsValid() {
			destValue.Field(i).Set(srcFieldValue)
		}
	}
}