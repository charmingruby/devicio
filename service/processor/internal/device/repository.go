package device

type RoutineRepository interface {
	Store(r *Routine) error
}
