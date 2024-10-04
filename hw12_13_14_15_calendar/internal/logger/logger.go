package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"github.com/TaoGunner/otus-hw/hw12_13_14_15_calendar/internal/config"
)

const (
	LoggerLevelDebug   string = "debug"
	LoggerLevelInfo    string = "info"
	LoggerLevelWarning string = "warning"
	LoggerLevelError   string = "error"

	LoggerFormatJSON   string = "json"
	LoggerFormatLOGFMT string = "logfmt"
)

var logLvl = new(slog.LevelVar)

func InitLogger(cfg config.Config) error {
	// Опции логгера
	opts := &slog.HandlerOptions{AddSource: cfg.Logger.IsAddSource, Level: logLvl, ReplaceAttr: replaceAttr}

	multiWriter := io.MultiWriter(os.Stdout)
	if len(cfg.Logger.LogFilename) > 0 {
		logFile, err := os.OpenFile(cfg.Logger.LogFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return err
		}

		multiWriter = io.MultiWriter(os.Stdout, logFile)
	}

	// Формат: JSON / logfmt
	switch cfg.Logger.Format {
	case LoggerFormatJSON:
		slog.SetDefault(slog.New(slog.NewJSONHandler(multiWriter, opts)))
	case LoggerFormatLOGFMT:
		slog.SetDefault(slog.New(slog.NewTextHandler(multiWriter, opts)))
	default:
		slog.SetDefault(slog.New(slog.NewTextHandler(multiWriter, opts)))
	}

	return nil
}

// SetLevel переключает уровень логирования.
func SetLevel(level string) {
	// Уровень логирования
	switch level {
	case LoggerLevelDebug:
		logLvl.Set(slog.LevelDebug)
	case LoggerLevelInfo:
		logLvl.Set(slog.LevelInfo)
	case LoggerLevelWarning:
		logLvl.Set(slog.LevelWarn)
	case LoggerLevelError:
		logLvl.Set(slog.LevelError)
	default:
		logLvl.Set(slog.LevelInfo)
	}
}

// replaceAttr заменяет стандартные аттрибуты логгера (time, source, level, msg).
func replaceAttr(_ []string, attr slog.Attr) slog.Attr {
	if attr.Key == slog.SourceKey {
		if sourceAttr, ok := attr.Value.Any().(*slog.Source); ok {
			_, f, l, _ := runtime.Caller(6)
			sourceAttr.Function = ""
			sourceAttr.File = filepath.Base(f)
			sourceAttr.Line = l
		}
	}

	return attr
}
