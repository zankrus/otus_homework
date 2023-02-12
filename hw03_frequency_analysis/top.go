package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

type Word struct {
	word  string
	count int
}

type WordSlice []Word

func Top10(source string) []string {
	wordsTank := make(WordSlice, 0)
	result := make([]string, 0)

	// Объект для счетчика слов
	wordsCount := map[string]int{}

	// Приведение к низкому регистру
	source = strings.ToLower(source)

	// Создаем массив слов
	words := strings.FieldsFunc(source, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != '-'
	})

	// Считаем количество вхождений слов
	for _, word := range words {
		if word == "-" {
			continue
		}

		if wordsCount[word] == 0 {
			wordsCount[word] = 1
			continue
		}
		wordsCount[word]++
	}

	// Создаем массив структур для удобства
	for word, count := range wordsCount {
		wordsTank = append(wordsTank, Word{
			word:  word,
			count: count,
		})
	}

	// Сортируем слайс по количеству слов или буквам
	sort.Slice(wordsTank, func(i, j int) bool {
		if wordsTank[j].count == wordsTank[i].count {
			return wordsTank[j].word > wordsTank[i].word
		}
		return wordsTank[j].count < wordsTank[i].count
	})

	// Обрезаем, если длинный массив
	if len(wordsTank) > 10 {
		wordsTank = wordsTank[:10]
	}

	for _, word := range wordsTank {
		result = append(result, word.word)
	}
	return result
}
