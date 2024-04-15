package utils

import (
	"fmt"
	"log"
	"reflect"
)

func InterfaceIsAStruct(i interface{}) bool {
	return reflect.ValueOf(i).Type().Kind() == reflect.Struct
}

func GetAllFieldOfAStruct(i interface{}) {
	if InterfaceIsAStruct(i) {
		return
	}
	val := reflect.ValueOf(i).Elem()
	for i := 0; i < val.NumField(); i++ {
		fmt.Println(val.Type().Field(i).Name)
		//fieldValue := val.Field(i)
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
