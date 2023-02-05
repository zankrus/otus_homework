package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
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

		if unicode.IsDigit(runedSourceString[i]) {

			number, err := strconv.Atoi(s)

			if err != nil {
				return "", ErrInvalidString
			}

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

		result += s
		isDigitLastChar = false

	}
	return result, nil
}
