package instrumentation

import (
	"log/slog"

	"github.com/charmingruby/devicio/lib/observability/log"
)

var (
	Logger *slog.Logger
)

func NewLogger() error {
	logger, err := log.NewSlogLogger("")
	if err != nil {
		return err
	}

	Logger = logger

	return nil
}
