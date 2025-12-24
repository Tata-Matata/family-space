package sqlite

import (
	"context"

	"database/sql"
	"errors"
	"time"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth"
	errs "github.com/Tata-Matata/family-space/apps/auth-service/internal/errors"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage"
)

type RefreshTokenStore struct {
	exec storage.SQLExecutor
}

func NewRefreshTokenStore(exec storage.SQLExecutor) storage.RefreshTokenStore {
	return &RefreshTokenStore{exec: exec}
}

func (s *RefreshTokenStore) Create(
	ctx context.Context,
	t auth.RefreshToken,
) error {

	query := `
		INSERT INTO refresh_tokens (
			id, user_id, token_hash,
			expires_at, revoked_at, created_at
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := s.exec.ExecContext(
		ctx,
		query,
		t.ID,
		t.UserID,
		t.TokenHash,
		t.ExpiresAt,
		t.RevokedAt,
		t.CreatedAt,
	)

	return err
}

func (s *RefreshTokenStore) GetByHash(
	ctx context.Context,
	hash string,
) (auth.RefreshToken, error) {

	query := `
		SELECT id, user_id, token_hash,
		       expires_at, revoked_at, created_at
		FROM refresh_tokens
		WHERE token_hash = ?
	`

	var t auth.RefreshToken
	var revoked sql.NullTime

	err := s.exec.QueryRowContext(ctx, query, hash).Scan(
		&t.ID,
		&t.UserID,
		&t.TokenHash,
		&t.ExpiresAt,
		&revoked,
		&t.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auth.RefreshToken{}, errs.ErrNotFound
		}
		return auth.RefreshToken{}, err
	}

	if revoked.Valid {
		t.RevokedAt = &revoked.Time
	}

	return t, nil
}

func (s *RefreshTokenStore) Revoke(
	ctx context.Context,
	id string,
) error {

	query := `
		UPDATE refresh_tokens
		SET revoked_at = ?
		WHERE id = ?
		  AND revoked_at IS NULL
	`

	res, err := s.exec.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrNotFound
	}

	return nil
}
