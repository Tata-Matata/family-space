package sqlite

import (
	"context"
	"database/sql"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
	errs "github.com/Tata-Matata/family-space/apps/auth-service/internal/errors"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage"
)

// the exact type of attached sql executor (sql.DB, sql.Tx etc)
// defines how the store will perform sql operations - in a transaction or not;
// this decision is made on the service layer (not on the store layer)
type UserStore struct {
	sql storage.SQLExecutor
}

func NewUserStore(sqlExec storage.SQLExecutor) storage.UserStore {
	return &UserStore{sql: sqlExec}
}

func (store *UserStore) Create(ctx context.Context, user domain.User) error {
	const q = `
		INSERT INTO users (id, email, password_hash, created_at)
		VALUES (?, ?, ?, ?)
	`

	_, err := store.sql.ExecContext(
		ctx,
		q,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}

func (store *UserStore) GetByEmail(ctx context.Context, email string) (domain.User, error) {

	const query = `
	  SELECT id, email, password_hash, created_at
	  FROM users
	  WHERE email = ?
	`

	var user domain.User
	err := store.sql.QueryRowContext(ctx, query, email).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, errs.ErrNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}

func (store *UserStore) GetById(ctx context.Context, id string) (domain.User, error) {

	const query = `
	  SELECT id, email, password_hash, created_at
	  FROM users
	  WHERE id = ?
	`

	var user domain.User
	err := store.sql.QueryRowContext(ctx, query, id).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, errs.ErrNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}
