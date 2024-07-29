package main

import (
	"log/slog"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		slog.Error("Для запуска требуется минимум 3 аргумента", "args_count", len(os.Args))
		os.Exit(1)
	}

	envs, err := ReadDir(os.Args[1])
	if err != nil {
		slog.Error("Ошибка чтения каталога", "path", os.Args[1])
		os.Exit(1)
	}

	os.Exit(RunCmd(os.Args[2:], envs))
}
