package sqlite

import (
	"context"
	"database/sql"
	"strings"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
	errs "github.com/Tata-Matata/family-space/apps/auth-service/internal/errors"
)

// We need SQLite stores that can operate inside a transaction.
// In contrast to the regular stores that use sql.DB, these stores receive a sql.Tx
type FamilyStoreTransaction struct {
	transaction *sql.Tx
}

func NewFamilyStoreTransaction(transaction *sql.Tx) *FamilyStoreTransaction {
	return &FamilyStoreTransaction{transaction: transaction}
}

func (store *FamilyStoreTransaction) Create(ctx context.Context, family domain.Family) error {
	const q = `
    INSERT INTO families (id, email, password_hash, created_at)
    VALUES (?, ?, ?, ?)
  `

	_, err := store.transaction.ExecContext(
		ctx,
		q,
		family.ID,
		family.Name,
		family.CreatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return errs.ErrAlreadyExists
		}
		return err
	}
	return nil
}
