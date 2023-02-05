package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inputString string) (string, error) {
	// Place your code here
	var result string
	var isDigitLastChar bool

	runedSourceString := []rune(inputString)

	for i := 0; i < len(runedSourceString); i++ {
		// Создаем переменную символа-строки
		s := string(runedSourceString[i])
		number, err := strconv.Atoi(s)

		if err != nil {
			// Если ошибка, значит символ не является числом
			result += s
			isDigitLastChar = false
		} else {
			// Если без ошибки, значит символ - цифра
			previousIndex := i - 1
			if previousIndex < 0 {
				return "", ErrInvalidString
			}
			if isDigitLastChar {
				return "", ErrInvalidString
			}
			if number == 0 {
				// Удаляем последний символ
				result = result[0 : i-1]
				continue
			}
			previousChar := string(runedSourceString[previousIndex])
			repeatedString := strings.Repeat(previousChar, number-1)
			result += repeatedString
			isDigitLastChar = true
			continue
		}
	}
	fmt.Println(result)
	return result, nil
}
