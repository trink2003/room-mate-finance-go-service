package utils

import (
	"fmt"
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
