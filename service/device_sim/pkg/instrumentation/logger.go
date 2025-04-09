package instrumentation

import (
	"github.com/charmingruby/devicio/lib/observability"
	"github.com/charmingruby/devicio/lib/observability/log"
)

var (
	Logger observability.Logger
)

func NewLogger(lvl string) {
	Logger = log.NewSlogLogger(lvl)
}
