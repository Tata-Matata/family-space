package sqlite

import (
	"context"
	"strings"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage"
)

// the exact type of attached sql executor (sql.DB, sql.Tx etc)
// defines how the store will perform sql operations - in a transaction or not;
// this decision is made on the service layer (not on the store layer)
type FamilyStore struct {
	sql storage.SQLExecutor
}

func NewFamilyStore(sqlExec storage.SQLExecutor) *FamilyStore {
	return &FamilyStore{sql: sqlExec}
}

func (store *FamilyStore) Create(ctx context.Context, family domain.Family) error {
	const query = `
	  INSERT INTO families (id, name, created_at)
	  VALUES (?, ?, ?, ?)
	`

	_, err := store.sql.ExecContext(
		ctx,
		query,
		family.ID,
		family.Name,
		family.CreatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return storage.ErrAlreadyExists
		}
		return err
	}

	return nil
}
