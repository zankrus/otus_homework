package main

import (
	"github.com/stretchr/testify/require" //nolint
	"log"
	"os"
	"testing" //nolint
)

func TestReadDir(t *testing.T) {
	t.Run("Директория с файлами", func(t *testing.T) {
		testDir := "./testdata/env"
		env, err := ReadDir(testDir)
		if err != nil {
			log.Fatal(err)
		}
		require.Equal(t, 5, len(env))
		require.Equal(t, EnvValue{
			Value:      "bar",
			NeedRemove: false,
		}, env["BAR"])
		require.Equal(t, EnvValue{
			Value:      "",
			NeedRemove: true,
		}, env["EMPTY"])
	})

	t.Run("Пустая директория", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "")
		if err != nil {
			t.Fail()
		}

		env, err := ReadDir(dir)

		require.NoError(t, err)
		require.Len(t, env, 0)
	})
	t.Run("Файл с =", func(t *testing.T) {
		file, err := os.CreateTemp("./testdata/env", "=")
		if err != nil {
			t.Fail()
		}

		_, err = ReadDir("./testdata/env")

		require.Error(t, err)

		os.RemoveAll(file.Name())
	})
}

func TestStringCleaner(t *testing.T) {
	const expected = "some"
	t.Run("Пробел вначале", func(t *testing.T) {
		testString := " some"
		result := stringCleaner(testString)
		require.Equal(t, expected, result)
	})
	t.Run("Пробел в конце", func(t *testing.T) {
		testString := "some "
		result := stringCleaner(testString)
		require.Equal(t, expected, result)
	})

	t.Run("Текст в две строки", func(t *testing.T) {
		testString := "some\n some again \n"
		result := stringCleaner(testString)
		require.Equal(t, expected, result)
	})

	t.Run("Текст с табуляцией ", func(t *testing.T) {
		testString := "some\t"
		result := stringCleaner(testString)
		require.Equal(t, expected, result)
	})
}
