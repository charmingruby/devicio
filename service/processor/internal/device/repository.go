package device

import "context"

type RoutineRepository interface {
	Store(ctx context.Context, r Routine) (context.Context, error)
}
