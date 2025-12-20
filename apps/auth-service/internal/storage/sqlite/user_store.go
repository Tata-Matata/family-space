package sqlite

import (
	"context"
	"database/sql"
	"strings"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (userStore *UserStore) Create(ctx context.Context, user domain.User) error {
	const query = `
	  INSERT INTO users (id, email, password_hash, created_at)
	  VALUES (?, ?, ?, ?)
	`

	_, err := userStore.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return storage.ErrAlreadyExists
		}
		return err
	}

	return nil
}

func (userStore *UserStore) GetByEmail(ctx context.Context, email string) (domain.User, error) {

	const query = `
	  SELECT id, email, password_hash, created_at
	  FROM users
	  WHERE email = ?
	`

	var user domain.User
	err := userStore.db.QueryRowContext(ctx, query, email).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, storage.ErrNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}
