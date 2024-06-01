package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

// Решение ДЗ №1 «Hello, OTUS!»
//
//	Необходимо написать программу, печатающую в стандартный вывод перевернутую фразу 'Hello, OTUS!'
func main() {
	fmt.Println(reverse.String("Hello, OTUS!"))
}
