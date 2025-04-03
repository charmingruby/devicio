package observability

import (
	"github.com/charmingruby/devicio/lib/observability"
	"github.com/charmingruby/devicio/lib/observability/otel"
)

var Tracer observability.Tracer

func NewTracer(serviceName string) error {
	tracer, err := otel.NewTracer(serviceName)
	if err != nil {
		return err
	}

	Tracer = tracer

	return nil
}
