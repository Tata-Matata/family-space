package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage/postgres"
)

func newTestFamily() domain.Family {
	return domain.Family{
		ID:        uuid.NewString(),
		Name:      "family-" + uuid.NewString(),
		CreatedAt: time.Now().UTC(),
	}
}

func TestFamilyStore_Create(test *testing.T) {
	db := newTestDB(test)
	store := postgres.NewFamilyStore(db)

	ctx := context.Background()
	family := newTestFamily()

	require.NoError(test, store.Create(ctx, family))
}

func TestFamilyStore_Create_AlreadyExists(test *testing.T) {
	db := newTestDB(test)
	store := postgres.NewFamilyStore(db)

	ctx := context.Background()
	family := newTestFamily()

	require.NoError(test, store.Create(ctx, family))
	err := store.Create(ctx, family)

	require.ErrorContains(test, err, "duplicate key value violates unique constraint")
}
