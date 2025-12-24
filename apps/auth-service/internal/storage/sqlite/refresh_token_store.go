package sqlite

import (
	"context"

	"database/sql"
	"errors"
	"time"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/refresh"
	errs "github.com/Tata-Matata/family-space/apps/auth-service/internal/errors"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage"
)

type RefreshTokenStore struct {
	exec storage.SQLExecutor
}

func NewRefreshTokenStore(exec storage.SQLExecutor) storage.RefreshTokenStore {
	return &RefreshTokenStore{exec: exec}
}

func (store *RefreshTokenStore) Create(
	ctx context.Context,
	token refresh.RefreshToken,
) error {

	query := `
		INSERT INTO refresh_tokens (
			id, user_id, token_hash,
			expires_at, revoked_at, created_at
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := store.exec.ExecContext(
		ctx,
		query,
		token.ID,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
		token.RevokedAt,
		token.CreatedAt,
	)

	return err
}

func (store *RefreshTokenStore) GetByHash(
	ctx context.Context,
	hash string,
) (refresh.RefreshToken, error) {

	query := `
		SELECT id, user_id, token_hash,
		       expires_at, revoked_at, created_at
		FROM refresh_tokens
		WHERE token_hash = ?
	`

	var token refresh.RefreshToken
	var revoked sql.NullTime

	err := store.exec.QueryRowContext(ctx, query, hash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&revoked,
		&token.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return refresh.RefreshToken{}, errs.ErrNotFound
		}
		return refresh.RefreshToken{}, err
	}

	if revoked.Valid {
		token.RevokedAt = &revoked.Time
	}

	return token, nil
}

func (store *RefreshTokenStore) Revoke(
	ctx context.Context,
	id string,
) error {

	query := `
		UPDATE refresh_tokens
		SET revoked_at = ?
		WHERE id = ?
		  AND revoked_at IS NULL
	`

	res, err := store.exec.ExecContext(ctx, query, time.Now(), id)
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
