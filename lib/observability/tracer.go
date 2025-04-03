package observability

import (
	"context"
)

type Tracer interface {
	Span(ctx context.Context, name string) (context.Context, func())
	Close() error
}
