package postgres

import (
	"context"
	"fmt"

	"github.com/charmingruby/devicio/lib/database"
	"github.com/charmingruby/devicio/lib/observability"
	"github.com/charmingruby/devicio/service/processor/internal/device"
	"github.com/charmingruby/devicio/service/processor/pkg/logger"
	"github.com/jmoiron/sqlx"
)

func NewRoutineRepository(db *sqlx.DB, tracer observability.Tracing) (*RoutineRepository, error) {
	stmts := make(map[string]*sqlx.Stmt)

	for queryName, statement := range routineQueries() {
		stmt, err := db.Preparex(statement)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("unable to prepare the query: %s, err: %s", queryName, err.Error()))
			return nil, database.ErrPreparation
		}

		stmts[queryName] = stmt
	}

	return &RoutineRepository{
		db:     db,
		tracer: tracer,
		stmts:  stmts,
	}, nil
}

type RoutineRepository struct {
	db     *sqlx.DB
	tracer observability.Tracing
	stmts  map[string]*sqlx.Stmt
}

func (r *RoutineRepository) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]

	if !ok {
		logger.Log.Error(fmt.Sprintf("statement not prepared: %s", queryName))
		return nil, database.ErrStatementNotPrepared
	}

	return stmt, nil
}

func (r *RoutineRepository) Store(ctx context.Context, routine device.Routine) (context.Context, error) {
	ctx, complete := r.tracer.Span(ctx, "repository.RoutineRepository.Store")
	defer complete()

	stmt, err := r.statement(createRoutine)
	if err != nil {
		return ctx, err
	}

	if _, err := stmt.Exec(
		routine.ID,
		routine.DeviceID,
		routine.Status,
		routine.Context,
		routine.Area,
		routine.DispatchedAt,
	); err != nil {
		return ctx, err
	}

	return ctx, nil
}
