package instrumentation

import (
	"log/slog"

	"github.com/charmingruby/devicio/lib/observability/log"
)

var (
	Logger *slog.Logger
)

func NewLogger(lvl string) {
	Logger = log.NewSlogLogger(lvl)
}
