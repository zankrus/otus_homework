package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {

	// Выбираем размер буффера
	bufferSize := 10
	buffer := make([]byte, bufferSize)

	// Открываем файл
	file, err := os.Open(fromPath)
	if err != nil {
		log.Fatal(err)
	}
	// Не забываем закрыть файл
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	// Устанавливаем отступ
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	// Создаем все папки для исходного файла
	err = os.MkdirAll(filepath.Dir(toPath), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// Создаем файл, в который будем копировать
	resFile, err := os.OpenFile(toPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	defer func(resFile *os.File) {
		err := resFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resFile)
	if err != nil {
		log.Fatal(err)
	}

	// счетчик буффера
	limitCount := int64(0)
	lastIteration := false

	for {
		limitCount += int64(bufferSize)
		if limitCount > limit {
			buffer = make([]byte, limitCount-limit)
			lastIteration = true
		}

		bytesRead, _ := file.Read(buffer)
		if bytesRead == 0 { // bytesRead будет равен 0 в конце файла.
			break
		}

		res := string(buffer)
		log.Println(res)

		// сразу писать в файл
		_, err := resFile.Write(buffer)
		if err != nil {
			log.Fatal(err)
		}

		if lastIteration {
			break
		}

	}

	return err
}
