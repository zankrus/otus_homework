package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
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
	errNotAStruct          = errors.New("переданный объект не структура")
	errMissmatchTagAndType = errors.New("тэг недопустим для типа")
	errBrokenTag           = errors.New("невалидный тэг")
	ErrMin                 = errors.New("значения меньше допустимого")
	ErrMax                 = errors.New("значения больше допустимого")
	ErrIn                  = errors.New("значение не входит в список допустимых")
	ErrLen                 = errors.New("допустима длина превышена")
	ErrRegex               = errors.New("строка не подходит под регулярное выражение")
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

		valErrors = checkErrors(field.Name, tags, fieldValue, valErrors)

		if len(valErrors) > 0 {
			return valErrors
		}
	}
	return nil
}

func checkErrors(fName string, fTags []string, fValue reflect.Value, errContainer []ValidationError) ValidationErrors {
	var errs []error
	var newValErrs = errContainer

	fmt.Println("***Функция checkErrors***")
	fmt.Printf("Имя поля  %v \n", fName)
	fmt.Printf("Теги поля  %v \n", fTags)
	fmt.Printf("Велью поля  %v \n", fValue)
	fmt.Printf("Тип поля  %v \n", fValue.Kind())
	fmt.Println("")

	switch fValue.Kind() {
	case reflect.Int:
		errs = validateByTag(fTags, fValue)
		fmt.Printf("Ошибки в поле  %v : %v \n", fName, errs)
	case reflect.String:
		errs = validateByTag(fTags, fValue)
		fmt.Printf("Ошибки в поле  %v : %v \n", fName, errs)
	case reflect.Slice:
		for i := 0; i < fValue.Len(); i++ {
			newValErrs = checkErrors(fName, fTags, fValue.Index(i), newValErrs)
		}

	}
	if len(errs) > 0 {
		for _, err := range errs {
			newValErrs = append(newValErrs, ValidationError{fName, err})
		}
	}

	return newValErrs

}

func validateByTag(tags []string, value reflect.Value) []error {
	var errs []error
	for _, tag := range tags {
		var err error
		splitedTag := strings.Split(tag, ":")
		if len(splitedTag) != 2 {
			err = errBrokenTag
			continue
		}
		tagName := splitedTag[0]
		tagValue := splitedTag[1]

		fmt.Println("***Функция validateByTag***")
		fmt.Printf("tagN %v \n", tagName)
		fmt.Printf("tagV %v \n", tagValue)

		switch tagName {

		case "min":
			err = compareInt(value, tagValue, "min")
		case "max":
			err = compareInt(value, tagValue, "max")
		case "len":
			err = lenCompare(value, tagValue)
		case "regexp":
			err = regexpCompare(value, tagValue)
		}
		if err != nil {
			errs = append(errs, err)
		}

	}
	return errs
}

func lenCompare(value reflect.Value, limit string) error {
	if value.Kind() != reflect.String {
		return errMissmatchTagAndType
	}
	limV, err := strconv.Atoi(limit)
	if err != nil {
		return err
	}
	if len(value.String()) > limV {
		return ErrLen
	}
	return nil
}
func regexpCompare(value reflect.Value, template string) error {
	if value.Kind() != reflect.String {
		return errMissmatchTagAndType
	}
	rx, err := regexp.Compile(template)
	if err != nil {
		return err
	}
	if !rx.MatchString(value.String()) {
		return ErrRegex
	}
	return nil
}

func compareInt(value reflect.Value, limit, operator string) error {
	if value.Kind() != reflect.Int {
		return errMissmatchTagAndType
	}
	limV, err := strconv.Atoi(limit)
	if err != nil {
		return err
	}

	switch operator {
	case "min":
		if int(value.Int()) < limV {
			return ErrMin
		}
	case "max":
		if int(value.Int()) > limV {
			return ErrMax
		}
	default:
		return errBrokenTag
	}
	return nil
}
