package application

import (
	"log/slog"
	"os"
)

func NewLogger(handler slog.Handler) *slog.Logger {
	return slog.New(handler)
}

func NewDefaultJsonLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
