package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type PostgresConnectionInput struct {
	User         string
	Password     string
	Host         string
	DatabaseName string
	SSL          string
}

func NewPostgres(in PostgresConnectionInput) (*sqlx.DB, error) {
	connectionString := fmt.Sprintf(
		"postgresql://%s:%s@%s/%s?sslmode=%s",
		in.User,
		in.Password,
		in.Host,
		in.DatabaseName,
		in.SSL,
	)
	dbDriver := "postgres"

	db, err := sqlx.Connect(dbDriver, connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
