package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

const validateTagKey = "validate"

type ValidationErrors []ValidationError

var (
	errNotAStruct = errors.New("переданный объект не структура")
)

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	objType := val.Type()

	//valErrors := make(ValidationErrors, 0)

	if objType.Kind() != reflect.Struct {
		return errNotAStruct
	}

	for i := 0; i < val.NumField(); i++ {

		field := objType.Field(i)
		fullTag, ok := field.Tag.Lookup(validateTagKey)

		//TODO: удалить принты отладочные
		fmt.Printf("field value is  %v ,  fulltag is %s \n", val.Field(i), fullTag)
		fmt.Printf("field is validating : %v \n", ok)
		fmt.Println("")

		if !ok {
			continue
		}

		tags := strings.Split(fullTag, "|")
		for _, tag := range tags {
			fmt.Println("")
			fmt.Printf("tag is : %v", tag)
		}
	}

	//TODO: удалить принты отладочные
	fmt.Println(val)
	fmt.Println(objType)

	return nil
}

//func intValidator(v int, expValues []int) error {
//
//}
