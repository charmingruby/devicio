package otel

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	tracer  trace.Tracer
	cleanup func() error
}

func NewTracer(serviceName string) (*Tracer, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	r := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
	)

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(r),
	)

	otel.SetTracerProvider(traceProvider)

	t := &Tracer{
		tracer: otel.Tracer(serviceName),
		cleanup: func() error {
			return traceProvider.Shutdown(context.Background())
		},
	}

	return t, nil
}

func (t *Tracer) Span(ctx context.Context, name string) (context.Context, func()) {
	ctx, span := t.tracer.Start(ctx, name)

	complete := func() {
		span.End()
	}

	return ctx, complete
}

func (t *Tracer) Close() error {
	return t.cleanup()
}
