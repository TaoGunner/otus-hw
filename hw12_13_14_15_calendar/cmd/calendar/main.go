package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TaoGunner/otus-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/TaoGunner/otus-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/TaoGunner/otus-hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/TaoGunner/otus-hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/TaoGunner/otus-hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/TaoGunner/otus-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/TaoGunner/otus-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.json", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	// Загрузка конфигурации сервиса
	cfg, err := config.NewConfig(configFile)
	if err != nil {
		slog.Error("Ошибка чтения файла конфигурации", "error", err)
		slog.Warn("Будут испоьзованы настройки по-умолчанию")
	}

	if err := logger.InitLogger(cfg); err != nil {
		slog.Error("ошибка инициализации логгера", "error", err)
	}

	var storage storage.EventStorer

	switch cfg.Database.Storage {
	case "in-memory":
		storage = memorystorage.New()
	case "sql":
		storage, err = sqlstorage.New(cfg.Database.DBPath)
		if err != nil {
			slog.Error("ошибка создания хранилища", "error", err)
			os.Exit(1)
		}
	default:
		slog.Error("неизвестный типа хранилища", "storage_type", cfg.Database.Storage)
		os.Exit(1)
	}

	calendar := app.New(storage)

	server := internalhttp.NewServer(cfg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			slog.Error("failed to stop http server", "error", err)
		}
	}()

	slog.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		slog.Error("failed to start http server", "error", err)
		cancel()
	}
}
