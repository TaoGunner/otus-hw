package main

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrInvalidString = errors.New("invalid string")
	ErrDigitParsing  = errors.New("digit parsing error")
)

// Решение ДЗ №2 «Распаковка строки»
//
//	Необходимо написать Go функцию, осуществляющую примитивную распаковку строки.
func Unpack(in string) (string, error) {
	var sb strings.Builder

	// Каст в руны для перебора через for (range не подходит)
	runes := []rune(in)
	for i := 0; i < len(runes); i++ {
		// Если проверяемый символ это цифра - это ошибка
		if unicode.IsDigit(runes[i]) {
			return "", ErrInvalidString
		}

		// Если за проверяемым символом идёт цифра
		if i < len(runes)-1 && unicode.IsDigit(runes[i+1]) {
			// Пытаемся распарсить руну в число (+ на всякий случай обработаем ошибку, хотя после IsDigit она маловерятна)
			repeatCount, err := strconv.Atoi(string(runes[i+1]))
			if err != nil {
				return "", ErrDigitParsing
			}

			// Продублируем руну repeatCount раз
			sb.WriteString(strings.Repeat(string(runes[i]), repeatCount))

			// Инкрементируем счетчик, потому что цифра за проверяемым символом нам уже не нужна
			i++

			continue
		}

		// В остальных случаях - просто сохраняем руну в результат
		sb.WriteRune(runes[i])
	}

	return sb.String(), nil
}
