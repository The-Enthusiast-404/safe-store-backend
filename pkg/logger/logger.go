package logger

import (
	"log/slog"
	"os"
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
}

type slogLogger struct {
	*slog.Logger
}

func New() Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return &slogLogger{slog.New(handler)}
}

func (l *slogLogger) Info(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}

func (l *slogLogger) Error(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}

func (l *slogLogger) Fatal(msg string, args ...any) {
	l.Logger.Info(msg, args...)
	os.Exit(1)
}
