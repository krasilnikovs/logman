package application

import (
	"log/slog"
	"os"
)

// NewLogger constructs new logger by passing slog handler
func NewLogger(handler slog.Handler) *slog.Logger {
	return slog.New(handler)
}

// NewDefaultJsonLogger constructs logger with json handler and output to stdout
func NewDefaultJsonLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
