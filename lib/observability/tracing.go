package observability

import (
	"context"
)

type Tracing interface {
	Span(ctx context.Context, name string) (context.Context, func())
	Close() error
}
