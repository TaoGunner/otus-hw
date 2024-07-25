package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	execCmd := exec.Command(cmd[0], cmd[1:]...)
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
		if exitError, ok := err.(*exec.ExitError); ok {
			returnCode = exitError.ExitCode()
		}
	}

	return
}

func modifyEnv(envSet Environment) []string {
	result := []string{}

	for _, osEnv := range os.Environ() {
		osEnvName := strings.SplitN(osEnv, "=", 2)[0]
		if env, ok := envSet[osEnvName]; ok {
			if !env.NeedRemove {
				result = append(result, fmt.Sprintf("%s=%s", osEnvName, env.Value))
			}
		} else {
			result = append(result, osEnv)
		}
	}

	return result
}
