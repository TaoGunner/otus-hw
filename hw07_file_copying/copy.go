package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// Размер батча для копирования.
const copyBatchSize int64 = 64

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrEqualPath             = errors.New("source and destination path are equal")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Пути файлов не должны быть одинаковыми
	if fromPath == toPath {
		return ErrEqualPath
	}

	// Исходный файл должен существовать
	srcFileInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	// Исходный файл не должен быть папкой и иметь ненулевой размер
	if srcFileInfo.IsDir() || srcFileInfo.Size() == 0 {
		return ErrUnsupportedFile
	}

	// Offset больше, чем размер файла - невалидная ситуация
	if offset > srcFileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	// Создаём будущий файл-копию
	dstFile, err := os.OpenFile(toPath, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Println(1)
		return err
	}
	defer dstFile.Close()

	// Открываем исходный файл
	srcFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0o644)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Переводим указатель на offset от начала файла
	if _, err := srcFile.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	// Посчитаем, сколько всего байт потребуется скопировать
	totalSize := srcFileInfo.Size() - offset
	if limit > 0 && limit < totalSize {
		totalSize = limit
	}

	copiedSize := int64(0)
	for copiedSize < totalSize {
		// Вычисление размера батча для копирования (по-умолчанию: <copyBatchSize> байт)
		batchSize := copyBatchSize
		if copiedSize+copyBatchSize > totalSize {
			batchSize = totalSize - copiedSize
		}

		// Копируем необходимое количество байт, игнорируя EOF потому что цикл и так завершится
		n, err := io.CopyN(dstFile, srcFile, batchSize)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}

		copiedSize += n

		// Рисуем шкалу прогресса
		drawProgressBar(copiedSize, totalSize)
	}

	return nil
}

// drawProgressBar отображает шкалу прогресса копирования в консоли.
func drawProgressBar(copiedSize, totalSize int64) {
	percent := copiedSize * 100 / totalSize
	fmt.Printf("\r[%-100s] %10d/%-10d %d%%", strings.Repeat("#", int(percent)), copiedSize, totalSize, percent)
	if copiedSize == totalSize {
		fmt.Print(" FINISHED\n")
	}
}
