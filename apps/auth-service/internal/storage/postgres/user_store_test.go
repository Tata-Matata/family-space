package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
	errs "github.com/Tata-Matata/family-space/apps/auth-service/internal/errors"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage/postgres"
)

func newTestUser() domain.User {
	return domain.User{
		ID:           uuid.NewString(),
		Email:        "user-" + uuid.NewString() + "@example.com",
		PasswordHash: "hash",
		CreatedAt:    time.Now().UTC(),
	}
}

func TestUserStore_CreateAndGetByEmail(test *testing.T) {
	db := newTestDB(test)
	store := postgres.NewUserStore(db)

	ctx := context.Background()
	user := newTestUser()

	require.NoError(test, store.Create(ctx, user))

	got, err := store.GetByEmail(ctx, user.Email)
	require.NoError(test, err)

	require.Equal(test, user.ID, got.ID)
	require.Equal(test, user.Email, got.Email)
}

func TestUserStore_GetByEmail_NotFound(test *testing.T) {
	db := newTestDB(test)
	store := postgres.NewUserStore(db)

	_, err := store.GetByEmail(context.Background(), "missing@example.com")
	require.ErrorIs(test, err, errs.ErrNotFound)
}

func TestUserStore_GetById_NotFound(test *testing.T) {
	db := newTestDB(test)
	store := postgres.NewUserStore(db)

	_, err := store.GetById(context.Background(), uuid.NewString())
	require.ErrorIs(test, err, errs.ErrNotFound)
}
