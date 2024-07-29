package main

import (
	"bytes"
	"os"
	"path"
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
	// Читаем каталог с переменными окружения
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// Создадим карту согласно количеству файлов
	envs := make(Environment, len(files))

	for _, file := range files {
		// Пропускаем каталоги и файлы с '=' в имени
		if file.IsDir() || strings.Contains(file.Name(), "=") {
			continue
		}

		// Читаем файл
		envValue, err := os.ReadFile(path.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		// Нормализуем значение переменной окружения
		envValue = normalize(envValue)

		// Добавляем её в карту
		envs[file.Name()] = EnvValue{
			Value:      string(envValue),
			NeedRemove: len(envValue) == 0,
		}
	}

	return envs, nil
}

func normalize(value []byte) []byte {
	value = bytes.Split(value, []byte{'\n'})[0]
	value = bytes.TrimRight(value, " \t")
	value = bytes.ReplaceAll(value, []byte{0x00}, []byte{'\n'})

	return value
}
