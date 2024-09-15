package logger

import (
	"testing"

	"github.com/TaoGunner/otus-hw/hw12_13_14_15_calendar/internal/config"
)

func TestLogger(t *testing.T) {
	t.Run("info logfmt without file", func(t *testing.T) {
		cfg := config.Config{
			Logger: config.CfgLogger{
				Level:       "INFO",
				Format:      "logfmt",
				IsAddSource: true,
			},
		}
		if err := InitLogger(cfg); err != nil {
			t.Error(err)
		}
	})

	t.Run("info logfmt with file", func(t *testing.T) {
		cfg := config.Config{
			Logger: config.CfgLogger{
				Level:       "INFO",
				Format:      "logfmt",
				IsAddSource: true,
				LogFilename: "application.json",
			},
		}
		if err := InitLogger(cfg); err != nil {
			t.Error(err)
		}
	})
}
