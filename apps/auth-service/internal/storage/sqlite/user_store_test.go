package sqlite

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage"
	"github.com/stretchr/testify/require"
)

func setupTestDB(test *testing.T) *UserStore {
	test.Helper()

	dbPath := "test_auth.db"
	os.Remove(dbPath)

	db, err := Open(dbPath)
	require.NoError(test, err)

	_, err = db.Exec(`
	  CREATE TABLE users (
		id TEXT PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL
	  );
	`)
	require.NoError(test, err)

	test.Cleanup(func() {
		db.Close()
		os.Remove(dbPath)
	})
	var _ storage.SQLExecutor = db
	var _ storage.SQLExecutor = (*sql.Tx)(nil)

	return NewUserStore(db)
}
func TestUserStore_CreateAndGet(test *testing.T) {
	store := setupTestDB(test)

	user := domain.User{
		ID:           "user-1",
		Email:        "anna@example.com",
		PasswordHash: "hash",
		CreatedAt:    time.Now(),
	}

	err := store.Create(context.Background(), user)
	require.NoError(test, err)

	loaded, err := store.GetByEmail(context.Background(), "anna@example.com")
	require.NoError(test, err)
	require.Equal(test, user.Email, loaded.Email)
}

func TestUserStore_DuplicateEmail(test *testing.T) {
	store := setupTestDB(test)

	user := domain.User{
		ID:           "user-1",
		Email:        "anna@example.com",
		PasswordHash: "hash",
		CreatedAt:    time.Now(),
	}

	require.NoError(test, store.Create(context.Background(), user))
	err := store.Create(context.Background(), user)

	require.ErrorIs(test, err, storage.ErrAlreadyExists)
}
