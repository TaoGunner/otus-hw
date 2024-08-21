package main

import (
	"errors"
	"io"
	"log/slog"
	"net"
	"time"
)

var errConnectionNotCreated = errors.New("connection not created")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnet struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnet{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
		conn:    nil,
	}
}

// Connect создаёт TCP-подключение по указанному адресу.
func (t *telnet) Connect() error {
	var err error
	if t.conn, err = net.DialTimeout("tcp", t.address, t.timeout); err != nil {
		return err
	}

	return nil
}

// Close разрывает TCP-подключение по указанному адресу.
func (t *telnet) Close() error {
	if t.conn == nil { // Проверка на созданное соединение
		return errConnectionNotCreated
	}

	if err := t.conn.Close(); err != nil {
		return err
	}

	return nil
}

// Send отправляет сообщение на указанный адрес.
func (t *telnet) Send() error {
	if t.conn == nil { // Проверка на созданное соединение
		return errConnectionNotCreated
	}

	if _, err := io.Copy(t.conn, t.in); err != nil {
		slog.Error("Send error", "error", err)

		return err
	}

	slog.Info("EOF")

	return nil
}

// Receive получает сообщение от указанного адреса.
func (t *telnet) Receive() error {
	if t.conn == nil { // Проверка на созданное соединение
		return errConnectionNotCreated
	}

	if _, err := io.Copy(t.out, t.conn); err != nil {
		slog.Error("Receive error", "error", err)

		return err
	}

	slog.Info("Connection was closed by peer")

	return nil
}
