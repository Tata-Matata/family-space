package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMembershipStore_CreateAndGet(test *testing.T) {
	db := newTestDB(test)
	store := postgres.NewMembershipStore(db)

	ctx := context.Background()
	m := domain.Membership{
		UserID:    uuid.NewString(),
		FamilyID:  uuid.NewString(),
		Role:      "owner",
		CreatedAt: time.Now().UTC(),
	}

	require.NoError(test, store.Create(ctx, m))

	got, err := store.GetByUserID(ctx, m.UserID)
	require.NoError(test, err)

	require.Equal(test, m.UserID, got.UserID)
	require.Equal(test, m.FamilyID, got.FamilyID)
	require.Equal(test, m.Role, got.Role)
}
