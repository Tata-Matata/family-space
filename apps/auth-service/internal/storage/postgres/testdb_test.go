package postgres_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

func newTestDB(test *testing.T) *sql.DB {
	test.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		test.Fatal("TEST_DATABASE_URL not set")
	}

	db, err := sql.Open("pgx", dsn)
	require.NoError(test, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	require.NoError(test, db.PingContext(ctx))

	// Clean DB before each test
	_, err = db.Exec(`
		TRUNCATE TABLE refresh_tokens RESTART IDENTITY CASCADE;
	`)
	require.NoError(test, err)

	test.Cleanup(func() {
		_ = db.Close()
	})

	return db
}
