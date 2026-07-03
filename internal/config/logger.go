package config

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

var Logger *slog.Logger
var TerminalLogger *slog.Logger

func InitLogger() {
	logDir := "storage/logs"
	if err := os.MkdirAll(logDir, 0o750); err != nil {
		slog.Error("Failed to create logs directory", "error", err)
		return
	}

	logFilePath := filepath.Join(logDir, "go.log")
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600) // #nosec G304 — path is hardcoded, no user input
	if err != nil {
		slog.Error("Failed to open log file", "error", err)
		return
	}

	level := slog.LevelInfo
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		level = slog.LevelDebug
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	opts := &slog.HandlerOptions{Level: level}
	multiWriter := io.MultiWriter(logFile, os.Stdout)

	Logger = slog.New(slog.NewTextHandler(multiWriter, opts))
	TerminalLogger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	slog.SetDefault(Logger)
}
