package database

import "errors"

var (
	ErrPreparation          = errors.New("unable to prepare the query")
	ErrStatementNotPrepared = errors.New("statement not prepared")
)
