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

	for i := 0; i < len(inputString); i++ {
		// Создаем переменную символа-строки
		s := string(inputString[i])

		// Пробуем перевести символ в цифру
		number, err := strconv.Atoi(s)

		if err != nil {
			//Если ошибка, значит символ не является числом
			result += s
			isDigitLastChar = false
		} else {
			//Если без ошибки, значит символ - цифра
			previousIndex := i - 1
			if previousIndex < 0 {
				return "", ErrInvalidString
			}
			if isDigitLastChar == true {
				return "", ErrInvalidString
			}
			if number == 0 {
				//Удаляем последний символ
				result = result[0 : i-1]
				continue
			}
			previousChar := string(inputString[previousIndex])
			repeatedString := strings.Repeat(previousChar, number-1)
			result += repeatedString
			isDigitLastChar = true
			continue
		}

	}
	fmt.Println(result)
	return result, nil
}
