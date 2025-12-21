package sqlite

import (
	"context"
	"database/sql"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage"
)

// the exact type of attached sql executor (sql.DB, sql.Tx etc)
// defines how the store will perform sql operations - in a transaction or not;
// this decision is made on the service layer (not on the store layer)
type MembershipStore struct {
	sql storage.SQLExecutor
}

func NewMembershipStore(sqlExec storage.SQLExecutor) *MembershipStore {
	return &MembershipStore{sql: sqlExec}
}

type Membership = domain.Membership

func (store *MembershipStore) Create(ctx context.Context, membership Membership) error {
	const q = `
		INSERT INTO memberships (user_id, family_id, role, created_at)
		VALUES (?, ?, ?, ?)
	`

	_, err := store.sql.ExecContext(
		ctx,
		q,
		membership.UserID,
		membership.FamilyID,
		membership.Role,
		membership.CreatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}

func (store *MembershipStore) GetByUserID(ctx context.Context, userID string) (Membership, error) {

	const query = `	
	  SELECT id, email, password_hash, created_at
	  FROM memberships
	  WHERE user_id = ?
	`

	var membership Membership
	err := store.sql.QueryRowContext(ctx, query, userID).
		Scan(&membership.UserID, &membership.FamilyID, &membership.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return Membership{}, storage.ErrNotFound
		}
		return Membership{}, err
	}

	return membership, nil
}
