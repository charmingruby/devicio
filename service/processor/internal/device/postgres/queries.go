package postgres

const (
	createRoutine = "create routine"
)

func routineQueries() map[string]string {
	return map[string]string{
		createRoutine: `INSERT INTO device_routines
		(id, device_id, status, context, area, dispatched_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING *`,
	}
}
