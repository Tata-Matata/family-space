package postgres_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/refresh"
	errs "github.com/Tata-Matata/family-space/apps/auth-service/internal/errors"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage/postgres"
)

func newTestToken() refresh.RefreshToken {
	now := time.Now().UTC()
	return refresh.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    uuid.NewString(),
		TokenHash: uuid.NewString(),
		ExpiresAt: now.Add(24 * time.Hour),
		CreatedAt: now,
		RevokedAt: nil,
	}
}

func TestRefreshTokenStore_CreateAndGetByHash(test *testing.T) {
	test.Logf("TEST_DATABASE_URL=%q", os.Getenv("TEST_DATABASE_URL"))
	db := newTestDB(test)
	store := postgres.NewRefreshTokenStore(db)

	ctx := context.Background()
	token := newTestToken()

	err := store.Create(ctx, token)
	require.NoError(test, err)

	got, err := store.GetByHash(ctx, token.TokenHash)
	require.NoError(test, err)

	require.Equal(test, token.ID, got.ID)
	require.Equal(test, token.UserID, got.UserID)
	require.Equal(test, token.TokenHash, got.TokenHash)
	require.WithinDuration(test, token.ExpiresAt, got.ExpiresAt, time.Second)
	require.Nil(test, got.RevokedAt)
}

func TestRefreshTokenStore_GetByHash_NotFound(test *testing.T) {
	db := newTestDB(test)
	store := postgres.NewRefreshTokenStore(db)

	_, err := store.GetByHash(context.Background(), "missing-hash")
	require.ErrorIs(test, err, errs.ErrNotFound)
}

func TestRefreshTokenStore_Revoke(test *testing.T) {
	db := newTestDB(test)
	store := postgres.NewRefreshTokenStore(db)

	ctx := context.Background()
	token := newTestToken()

	require.NoError(test, store.Create(ctx, token))

	err := store.Revoke(ctx, token.ID)
	require.NoError(test, err)

	got, err := store.GetByHash(ctx, token.TokenHash)
	require.NoError(test, err)
	require.NotNil(test, got.RevokedAt)
}

func TestRefreshTokenStore_Revoke_NotFound(test *testing.T) {
	db := newTestDB(test)
	store := postgres.NewRefreshTokenStore(db)

	err := store.Revoke(context.Background(), uuid.NewString())
	require.ErrorIs(test, err, errs.ErrNotFound)
}
