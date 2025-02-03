package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type (
	LogFormat string
	LogLevel  string
)

const (
	TextFormat   LogFormat = "text"
	JSONFormat   LogFormat = "json"
	StderrOutput string    = "@stderr"
	StdoutOutput string    = "@stdout"
	DebugLevel   LogLevel  = "debug"
	InfoLevel    LogLevel  = "info"
	WarnLevel    LogLevel  = "warn"
	ErrorLevel   LogLevel  = "error"
)

var levels = map[LogLevel]int{
	DebugLevel: 0,
	InfoLevel:  1,
	WarnLevel:  2,
	ErrorLevel: 3,
}

type Logger struct {
	format LogFormat
	level  LogLevel
	writer io.Writer
	file   *os.File
}

type LogRecord struct {
	Timestamp   time.Time `json:"timestamp"`
	Level       string    `json:"level"`
	ClientIP    string    `json:"client_ip,omitempty"`
	Method      string    `json:"method,omitempty"`
	Path        string    `json:"path,omitempty"`
	HTTPVersion string    `json:"http_version,omitempty"`
	StatusCode  int       `json:"status_code,omitempty"`
	Latency     float64   `json:"latency_ms,omitempty"`
	UserAgent   string    `json:"user_agent,omitempty"`
	Message     string    `json:"message,omitempty"`
}

func NewLogger(outputPath string, level string, format LogFormat) (*Logger, error) {
	var writer io.Writer
	var file *os.File

	// Determine output destination.
	switch outputPath {
	case StderrOutput:
		writer = os.Stderr
	case StdoutOutput:
		writer = os.Stdout
	default:
		var err error
		file, err = os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		writer = file
	}

	return &Logger{
		writer: writer,
		format: format,
		level:  LogLevel(strings.ToLower(level)),
		file:   file,
	}, nil
}

func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

func (l *Logger) write(record LogRecord) error {
	if l.format == JSONFormat {
		encoder := json.NewEncoder(l.writer)
		return encoder.Encode(record)
	}

	// Text format.
	var msg string
	if record.ClientIP != "" {
		// HTTP request log.
		msg = fmt.Sprintf("%s | %s | %s | %s %s %s | %d | %.2fms | %s | %s\n",
			record.Timestamp.Format(time.RFC3339),
			record.Level,
			record.ClientIP,
			record.Method,
			record.Path,
			record.HTTPVersion,
			record.StatusCode,
			record.Latency,
			record.UserAgent,
			record.Message,
		)
	} else {
		// Simple application log.
		msg = fmt.Sprintf("%s | %s | %s\n",
			record.Timestamp.Format(time.RFC3339),
			record.Level,
			record.Message,
		)
	}

	_, err := fmt.Fprint(l.writer, msg)
	return err
}

// shouldLog determines if the message should be logged based on level.
func (l *Logger) shouldLog(msgLevel LogLevel) bool {
	return levels[msgLevel] >= levels[l.level]
}

func (l *Logger) log(level LogLevel, msg string) {
	if !l.shouldLog(level) {
		return
	}

	record := LogRecord{
		Timestamp: time.Now(),
		Level:     string(level),
		Message:   msg,
	}

	if err := l.write(record); err != nil {
		// If we can't log to our configured writer, fall back to stderr.
		fmt.Fprintf(os.Stderr, "Failed to write log: %v\n", err)
	}
}

// Debug logs a debug message.
func (l *Logger) Debug(msg string) {
	l.log(DebugLevel, msg)
}

// Info logs an info message.
func (l *Logger) Info(msg string) {
	l.log(InfoLevel, msg)
}

// Warn logs a warning message.
func (l *Logger) Warn(msg string) {
	l.log(WarnLevel, msg)
}

// Error logs an error message.
func (l *Logger) Error(msg string) {
	l.log(ErrorLevel, msg)
}

// LogRequest logs an HTTP request.
func (l *Logger) LogRequest(record LogRecord) {
	if !l.shouldLog(LogLevel(strings.ToLower(record.Level))) {
		return
	}

	if err := l.write(record); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write request log: %v\n", err)
	}
}
