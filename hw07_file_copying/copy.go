package main

import (
	"errors"
	"log"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {

	// Выбираем размер буффера
	bufferSize := 10
	b := make([]byte, bufferSize)

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

	// счетчик буффера

	for {
		bytesRead, _ := file.Read(b)
		if bytesRead == 0 { // bytesRead будет равен 0 в конце файла.
			break
		}

		// file[offset: limit] python

		res := string(b)
		log.Println(res)
		// сразу писать в файл

	}

	return err
}
