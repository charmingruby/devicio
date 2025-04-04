package observability

import (
	"context"
	"log/slog"
)

type Logger = slog.Logger

type Tracer interface {
	Span(ctx context.Context, name string) (context.Context, func())
	GetTraceIDFromContext(ctx context.Context) string
	Close() error
}

type Meter interface{}

type Instrumentation struct {
	Tracer Tracer
	Logger Logger
	Meter  Meter
}
