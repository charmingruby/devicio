package observability

import "context"

type Tracer interface {
	Span(ctx context.Context, name string) (context.Context, func())
	GetTraceIDFromContext(ctx context.Context) string
	Close() error
}

type Instrumentation struct {
	Tracer Tracer
}
