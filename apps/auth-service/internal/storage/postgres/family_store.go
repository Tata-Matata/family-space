package postgres

import (
	"context"

	"github.com/jackc/pgconn"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
	errs "github.com/Tata-Matata/family-space/apps/auth-service/internal/errors"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage"
)

// the exact type of attached sql executor (sql.DB, sql.Tx etc)
// defines how the store will perform sql operations - in a transaction or not;
// this decision is made on the service layer (not on the store layer)
type FamilyStore struct {
	sql storage.SQLExecutor
}

func NewFamilyStore(sqlExec storage.SQLExecutor) storage.FamilyStore {
	return &FamilyStore{sql: sqlExec}
}

func (store *FamilyStore) Create(
	ctx context.Context,
	family domain.Family,
) error {

	const query = `
		INSERT INTO families (
			id,
			name,
			created_at
		)
		VALUES ($1, $2, $3)
	`

	_, err := store.sql.ExecContext(
		ctx,
		query,
		family.ID,
		family.Name,
		family.CreatedAt,
	)
	if err != nil {
		// Postgres unique violation
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return errs.ErrAlreadyExists
			}
		}
		return err
	}

	return nil
}
