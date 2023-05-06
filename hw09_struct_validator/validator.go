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

func (v ValidationError) Error() string {
	return fmt.Sprintf("%v ошибка в поле %s", v.Err, v.Field)
}

const validateTagKey = "validate"

type ValidationErrors []ValidationError

var (
	errNotAStruct = errors.New("переданный объект не структура")
)

var (
	ErrLen            = errors.New("length")
	ErrRegex          = errors.New("regex")
	ErrMin            = errors.New("greater")
	ErrMax            = errors.New("less")
	ErrIn             = errors.New("lots of")
	ErrInvalidValues  = errors.New("invalid values")
	ErrExpectedStruct = errors.New("expected a struct")
)

func (v ValidationErrors) Error() string {
	b := strings.Builder{}
	for _, e := range v {
		b.WriteString(e.Err.Error())
		b.WriteRune('\n')
	}
	return b.String()
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	objType := val.Type()

	valErrors := make(ValidationErrors, 0)

	if objType.Kind() != reflect.Struct {
		return errNotAStruct
	}

	for i := 0; i < val.NumField(); i++ {
		fieldValue := val.Field(i)
		field := objType.Field(i)
		fullTag, ok := field.Tag.Lookup(validateTagKey)

		if !ok {
			continue
		}

		tags := strings.Split(fullTag, "|")

		valErrors = checkErrors(field.Name, tags, fieldValue)

		if len(valErrors) > 0 {
			return valErrors
		}
	}
	return nil
}

func checkErrors(fName string, fTags []string, fValue reflect.Value) ValidationErrors {

	fmt.Println("***Функция checkErrors***")
	fmt.Printf("Имя поля  %v \n", fName)
	fmt.Printf("Теги поля  %v \n", fTags)
	fmt.Printf("Велью поля  %v \n", fValue)
	fmt.Printf("Тип поля  %v \n", fValue.Kind())
	fmt.Println("")

	switch fValue.Kind() {
	case reflect.Int:

	case reflect.String:

	case reflect.Slice:

	}

	return nil
}
