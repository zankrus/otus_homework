package main

import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := Environment{}

	// Читаем директорию с файлами
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: проверка на пустую директорию

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Читаем файл
		fileName := file.Name()

		// Проверка на символ `=` в названии файла
		if strings.Contains(fileName, "=") {
			return nil, errors.New("File name \"" + fileName + "\" contains \"=\"")
		}
		file, err := os.Open(dir + "/" + file.Name())
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(file)
		scanner.Scan()

		content := bytes.Replace(scanner.Bytes(), []byte{0x00}, []byte("\n"), 1)
		value := string(bytes.TrimRight(content, " \t\r"))

		// Чистим файл от мусора
		if value == "" {
			env[fileName] = EnvValue{
				Value:      "",
				NeedRemove: true,
			}
			continue
		}

		env[fileName] = EnvValue{
			Value:      value,
			NeedRemove: false,
		}
	}

	return env, nil
}

func stringCleaner(s string) string {
	value := strings.Split(s, "\n")[0]
	value = strings.TrimLeft(value, " ")
	value = strings.TrimRight(value, "\t")
	value = strings.TrimRight(value, "\r")
	value = strings.TrimRight(value, " ")
	return value
}
