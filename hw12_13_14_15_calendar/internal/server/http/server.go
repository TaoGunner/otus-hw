package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/TaoGunner/otus-hw/hw12_13_14_15_calendar/internal/config"
)

const (
	restTimeoutDuration = time.Duration(10) * time.Second
)

type Server struct {
	http.Server
}

type Application interface { // TODO
}

func NewServer(cfg config.Config, _ Application) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handlerHello)

	handler := loggingMiddleware(mux)

	s := &Server{
		Server: http.Server{
			Addr:         fmt.Sprintf("%s:%d", cfg.Rest.Host, cfg.Rest.Port),
			Handler:      handler,
			ReadTimeout:  restTimeoutDuration,
			WriteTimeout: restTimeoutDuration,
			IdleTimeout:  restTimeoutDuration,
		},
	}

	return s
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.ListenAndServe(); err != nil {
		// Обработка ситуации, если порт REST-сервера занят
		var opErr *net.OpError
		var sysCallErr *os.SyscallError
		if errors.As(err, &opErr) && errors.As(opErr.Err, &sysCallErr) {
			if errors.Is(sysCallErr.Err, syscall.EADDRINUSE) {
				slog.Error("Starting HTTP server error (port is busy)")
			}
		} else {
			slog.Error("HTTP server error", "error", err)
		}
	}

	<-ctx.Done()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.Shutdown(ctx); err != nil {
		slog.Error("Stopping HTTP server error", "error", err)
	}

	return nil
}

func handlerHello(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello, world!"))
}
