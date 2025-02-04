package internalhttp

import (
	"net/http"
	"time"

	//nolint:depguard
	"github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/logger"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// Middleware for HTTP request logging.
func LoggingMiddleware(l *logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Create a custom response writer to capture the status code
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Process the request
		next.ServeHTTP(rw, r)

		// Calculate request duration
		duration := time.Since(startTime)

		// Create log record
		record := logger.LogRecord{
			Timestamp:   startTime,
			Level:       string(logger.InfoLevel),
			ClientIP:    r.RemoteAddr,
			Method:      r.Method,
			Path:        r.URL.Path,
			HTTPVersion: r.Proto,
			StatusCode:  rw.statusCode,
			Latency:     float64(duration.Milliseconds()),
			UserAgent:   r.UserAgent(),
		}

		// Log the record
		l.LogRequest(record)
	})
}
