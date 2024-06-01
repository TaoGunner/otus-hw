package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

// Решение ДЗ №1 «Hello, OTUS!»
//
//	Необходимо написать программу, печатающую в стандартный вывод перевернутую фразу 'Hello, OTUS!'.
//
//	2024-06-01
//	 - пакет stringutil переехал, но чек-лист и страница 'https://otus.ru/learning/301960' еще не обновлены.
//	 - можно использовать print() вместо fmt.Println, но это ломает ./test.sh
func main() {
	fmt.Println(reverse.String("Hello, OTUS!"))
}
