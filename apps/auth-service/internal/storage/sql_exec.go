package storage

import (
	"context"
	"database/sql"
)

// Both *sql.DB and *sql.Tx (transactional) implement this interface
// allowing us to write stores with either of them attached as behavior
type SQLExecutor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}
