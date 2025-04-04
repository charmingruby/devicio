package observability

import (
	"context"
)

type Tracing interface {
	Span(ctx context.Context, name string) (context.Context, func())
	GetTraceIDFromContext(ctx context.Context) string
	Close() error
}
