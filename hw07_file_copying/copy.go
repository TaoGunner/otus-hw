package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Размер батча для копирования.
const copyBatchSize int64 = 64

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
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

	// Сравним пути источника и приемника
	isEqual, err := filepath.Match(fromPath, toPath)
	if err != nil {
		return err
	}
	// Если равны - добавим к приемнику приставку '.tmp'
	if isEqual {
		toPath += ".tmp"
	}

	// Создаём будущий файл-копию
	dstFile, err := os.OpenFile(toPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
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
	percent := 0
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

		// Рисуем шкалу прогресса, если она изменилась хотя бы на 1%
		if int(copiedSize*100/totalSize) != percent {
			percent = int(copiedSize * 100 / totalSize)
			fmt.Printf("\r[%-100s] %d%%", strings.Repeat("#", percent), percent)
		}
	}

	// Если источник и приемник одинаковы
	if isEqual {
		// Удалим источник
		if err := os.Remove(fromPath); err != nil {
			return err
		}
		// Переименуем приемник
		if err := os.Rename(toPath, fromPath); err != nil {
			return err
		}
	}

	return nil
}
