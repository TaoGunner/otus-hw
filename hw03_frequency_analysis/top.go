package hw03frequencyanalysis

import (
	"regexp"
	"sort"

	"golang.org/x/exp/maps"
)

// Регулярное выражение собирающее слова, состоящие из видимых символов.
var regexpSimpleWords = regexp.MustCompile(`(\S+)`)

// Решение ДЗ №3 «Частотный анализ»
//
//	Необходимо написать Go функцию, принимающую на вход строку с текстом и
//	возвращающую слайс с 10-ю наиболее часто встречаемыми в тексте словами.
//	Если слова имеют одинаковую частоту, то должны быть отсортированы лексикографически.
func Top10(in string) []string {
	// Подсчёт и заполнение карты map[<уникальное слово>]<кол-во вхождений>
	wordsSet := map[string]int{}
	for _, word := range regexpSimpleWords.FindAllString(in, -1) {
		wordsSet[word]++
	}

	// Заполнение карты map[<кол-во вхождений>]<массив уникальных слов>
	topSet := map[int][]string{}
	for word, count := range wordsSet {
		// Создание новой связки ключ-значение, если она отсутствует в карте
		if _, ok := topSet[count]; !ok {
			topSet[count] = []string{}
		}

		// Добавление слова в массив встречающихся с такой же частотой
		words := topSet[count]
		words = append(words, word)
		topSet[count] = words
	}

	// Обратная сортировка массива ключей карты (пример: [9 7 6 5 3 2 1])
	countList := maps.Keys(topSet)
	sort.Sort(sort.Reverse(sort.IntSlice(countList)))

	// Заполнение результирующего массива 'result' лексикографически отсортированными словами
	result := []string{}
	for _, count := range countList {
		sort.Strings(topSet[count]) // Лексикографическая сортировка
		result = append(result, topSet[count]...)

		// Возврат среза первых 10 элементов результирующего массива, если его размер это позволяет
		if len(result) >= 10 {
			return result[:10]
		}
	}

	return result
}
