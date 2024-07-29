package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	execCmd := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	execCmd.Env = os.Environ()
	// Стандартные потоки ввода/вывода/ошибок пробрасываем в вызываемую программу
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	// Добавление (при необходимости - затирание) переменных окружения
	for key, value := range env {
		execCmd.Env = append(execCmd.Env, fmt.Sprintf("%s=%s", key, value.Value))
	}

	if err := execCmd.Run(); err != nil {
		// Код выхода утилиты приравниваем с кодом выхода программы
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			returnCode = exitError.ExitCode()
		}
	}

	return
}
