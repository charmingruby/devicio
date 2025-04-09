package postgres

import (
	"context"
	"fmt"

	"github.com/charmingruby/devicio/lib/database"
	"github.com/charmingruby/devicio/service/processor/internal/device"
	"github.com/charmingruby/devicio/service/processor/pkg/instrumentation"
	"github.com/jmoiron/sqlx"
)

func NewRoutineRepository(db *sqlx.DB) (*RoutineRepository, error) {
	stmts := make(map[string]*sqlx.Stmt)

	for queryName, statement := range routineQueries() {
		stmt, err := db.Preparex(statement)
		if err != nil {
			instrumentation.Logger.Error(fmt.Sprintf("unable to prepare the query: %s, err: %s", queryName, err.Error()))
			return nil, database.ErrPreparation
		}

		stmts[queryName] = stmt
	}

	return &RoutineRepository{
		db:    db,
		stmts: stmts,
	}, nil
}

type RoutineRepository struct {
	db    *sqlx.DB
	stmts map[string]*sqlx.Stmt
}

func (r *RoutineRepository) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]

	if !ok {
		instrumentation.Logger.Error(fmt.Sprintf("statement not prepared: %s", queryName))
		return nil, database.ErrStatementNotPrepared
	}

	return stmt, nil
}

func (r *RoutineRepository) Store(ctx context.Context, routine device.Routine) (context.Context, error) {
	ctx, complete := instrumentation.Tracer.Span(ctx, "repository.RoutineRepository.Store")
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
