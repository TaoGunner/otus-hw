package main

import (
	"context"
	"flag"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	flagTimeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	if len(flag.Args()) != 2 {
		slog.Error("invalid arguments count")

		os.Exit(1)
	}

	address := net.JoinHostPort(flag.Arg(0), flag.Arg(1))
	telnetClient := NewTelnetClient(address, *flagTimeout, os.Stdin, os.Stdout)

	if err := telnetClient.Connect(); err != nil {
		slog.Error("Connection error", "error", err)

		os.Exit(1)
	}
	defer telnetClient.Close()

	slog.Info("Connected to", "address", address)

	// При получении `SIGINT` программа должна завершать свою работу.
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()

	go func() {
		telnetClient.Send()
		cancel()
	}()

	go func() {
		telnetClient.Receive()
		cancel()
	}()

	// Ожидание сигнала SIGINT
	<-ctx.Done()
}
