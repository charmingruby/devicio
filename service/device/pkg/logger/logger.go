package logger

import (
	"log/slog"
	"os"
)

var (
	Log *slog.Logger
)

func New() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	Log = logger
}
