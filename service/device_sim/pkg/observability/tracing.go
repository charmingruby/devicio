package observability

import (
	"github.com/charmingruby/devicio/lib/observability"
	"github.com/charmingruby/devicio/lib/observability/otel"
)

var Tracing observability.Tracing

func NewTracing(serviceName string) error {
	tracing, err := otel.NewTracing(serviceName)
	if err != nil {
		return err
	}

	Tracing = tracing

	return nil
}
