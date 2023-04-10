package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrEmptyPaths            = errors.New("empty path")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" || toPath == "" {
		return ErrEmptyPaths
	}

	// Выбираем размер буффера
	bufferSize := int64(10)
	if limit != 0 && limit < bufferSize {
		bufferSize = limit
	}

	// Открываем файл
	file, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	// Не забываем закрыть файл
	defer file.Close()

	// Проводим валидации
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	if !fileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if fileInfo.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	// Устанавливаем отступ
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	// Создаем все папки для исходного файла
	err = os.MkdirAll(filepath.Dir(toPath), os.ModePerm)
	if err != nil {
		return err
	}

	// Создаем файл, в который будем копировать
	resFile, err := os.OpenFile(toPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	defer resFile.Close() //nolint:staticcheck

	if err != nil {
		return err
	}

	// счетчик буффера
	limitCount := int64(0)

	// Прогресс бар инит
	var expectedLen int64
	expectedLen = fileInfo.Size() - offset

	if expectedLen > limit && limit != 0 {
		expectedLen = limit
	}

	bar := pb.StartNew(int(expectedLen))
	bar.SetRefreshRate(time.Nanosecond)
	defer bar.Finish()

	for {
		// сразу писать в файл
		written, err := io.CopyN(resFile, file, bufferSize)

		bar.Add(int(written))

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		limitCount += written
		if limitCount == limit {
			break
		}
	}

	bar.Finish()
	return err
}
