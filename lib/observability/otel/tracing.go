package otel

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type Tracing struct {
	tracer  trace.Tracer
	cleanup func() error
}

func NewTracing(serviceName string) (*Tracing, error) {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint())
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

	t := &Tracing{
		tracer: otel.Tracer(serviceName),
		cleanup: func() error {
			return traceProvider.Shutdown(context.Background())
		},
	}

	return t, nil
}

func (t *Tracing) Span(ctx context.Context, name string) (context.Context, func()) {
	ctx, span := t.tracer.Start(ctx, name)

	complete := func() {
		span.End()
	}

	return ctx, complete
}

func (t *Tracing) GetTraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return ""
	}

	return span.SpanContext().TraceID().String()
}

func (t *Tracing) Close() error {
	return t.cleanup()
}
