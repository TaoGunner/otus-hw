package internalhttp

import (
	"log/slog"
	"net"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestTime := time.Now()

		// Получение IP-адреса пользователя
		ipAddr, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			slog.Error("Getting client ip error", "address", r.RemoteAddr, "error", err)
		}

		wrec := statusRecorder{w, 200}
		next.ServeHTTP(&wrec, r)

		slog.Info("API-Request",
			"ip", ipAddr,
			"method", r.Method,
			"uri", r.RequestURI,
			"user-agent", r.UserAgent(),
			"protcol", r.Proto,
			"status", wrec.status,
			"duration", time.Since(requestTime).Milliseconds(),
		)
	})
}

// Перехватчик статуса.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}
