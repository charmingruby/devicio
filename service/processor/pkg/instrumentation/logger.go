package instrumentation

import (
	"log/slog"
	"os"
)

var (
	Logger *slog.Logger
)

func NewLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	Logger = logger
}
