package instrumentation

import (
	"log/slog"

	"github.com/charmingruby/devicio/lib/observability/log"
)

var (
	Logger *slog.Logger
)

func NewLogger(lvl string) error {
	logger, err := log.NewSlogLogger(lvl)
	if err != nil {
		return err
	}

	Logger = logger

	return nil
}
