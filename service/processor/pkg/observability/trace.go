package observability

import (
	"github.com/charmingruby/devicio/lib/observability"
	"github.com/charmingruby/devicio/lib/observability/trace"
)

var Tracer observability.Tracer

func NewTracer(serviceName string) error {
	tracer, err := trace.NewOtelTracer(serviceName)
	if err != nil {
		return err
	}

	Tracer = tracer

	return nil
}
