package utils

import (
	"fmt"
	"log"
	"reflect"
)

func InterfaceIsAStruct(i interface{}) bool {
	return reflect.ValueOf(i).Type().Kind() == reflect.Struct
}

func GetAllFieldOfAStruct(inputInterface interface{}) {
	if !InterfaceIsAStruct(inputInterface) {
		return
	}
	reflect.TypeOf(inputInterface)
	val := reflect.ValueOf(inputInterface) /*.Elem()*/
	for i := 0; i < val.NumField(); i++ {
		currentFieldStructField := val.Type().Field(i)
		currentFieldValue := val.Field(i)
		fmt.Printf(
			"field information:\n\t>> field name: %s\n\t>> field value: %s\n\t>> field tag: %s\n\t>> example value in field tag json: %s\n",
			currentFieldStructField.Name,
			currentFieldValue,
			currentFieldStructField.Tag,
			GetStructTagValue(currentFieldStructField, "json"),
		)
	}

}

func SetValue(obj any, field string, value any) {
	ref := reflect.ValueOf(obj)

	// if its a pointer, resolve its value
	if ref.Kind() == reflect.Ptr {
		ref = reflect.Indirect(ref)
	}

	if ref.Kind() == reflect.Interface {
		ref = ref.Elem()
	}

	// should double-check we now have a struct (could still be anything)
	if ref.Kind() != reflect.Struct {
		log.Printf("unexpected type")
		return
	}

	prop := ref.FieldByName(field)
	prop.Set(reflect.ValueOf(value))
}

func GetStructTagValue(f reflect.StructField, tagName string) string {
	return f.Tag.Get(tagName)
}
