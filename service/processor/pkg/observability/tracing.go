package observability

import (
	"github.com/charmingruby/devicio/lib/observability"
	"github.com/charmingruby/devicio/lib/observability/otel"
)

var Tracing observability.Tracing

func NewTracing(serviceName string) error {
	tracer, err := otel.NewTracing(serviceName)
	if err != nil {
		return err
	}

	Tracing = tracer

	return nil
}
