package postgres

import (
	"context"
	"fmt"

	"github.com/charmingruby/devicio/service/processor/internal/device"
	"github.com/charmingruby/devicio/service/processor/pkg/logger"
	"github.com/charmingruby/devicio/service/processor/pkg/observability"
	"github.com/jmoiron/sqlx"
)

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

func NewRoutineRepository(db *sqlx.DB) (*RoutineRepository, error) {
	stmts := make(map[string]*sqlx.Stmt)

	for queryName, statement := range routineQueries() {
		stmt, err := db.Preparex(statement)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("unable to prepare the query: %s, err: %s", queryName, err.Error()))
			return nil, ErrPreparation
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
		logger.Log.Error(fmt.Sprintf("statement not prepared: %s", queryName))
		return nil, ErrStatementNotPrepared
	}

	return stmt, nil
}

func (r *RoutineRepository) Store(ctx context.Context, routine *device.Routine) (context.Context, error) {
	ctx, span := observability.Tracer.Start(ctx, "repository.RoutineRepository.Store")
	defer span.End()

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
