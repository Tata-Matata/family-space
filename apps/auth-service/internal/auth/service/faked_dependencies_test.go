package service_test

import (
	"context"
	"database/sql"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
	errs "github.com/Tata-Matata/family-space/apps/auth-service/internal/errors"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage"
)

type User = domain.User
type Membership = domain.Membership

/**FAKE DATABASE FOR UNIT TESTS **/
type fakeSQLExecutor struct{}

func (f *fakeSQLExecutor) ExecContext(
	ctx context.Context,
	query string,
	args ...any,
) (sql.Result, error) {
	panic("ExecContext should not be called in service unit test")
}

func (f *fakeSQLExecutor) QueryRowContext(
	ctx context.Context,
	query string,
	args ...any,
) *sql.Row {
	panic("QueryRowContext should not be called in service unit test")
}

type fakeDB struct {
	exec storage.SQLExecutor
}

// immitates starting db transaction but actually does nothing
func (f *fakeDB) BeginTransaction(
	ctx context.Context,
	readOnly bool,
) (storage.SQLExecutor, func(error) error, error) {
	return f.exec, func(error) error { return nil }, nil
}

/********** USER STORE **********/
type fakeUserStore struct {
	user User
	err  error
}

func (f *fakeUserStore) GetByEmail(ctx context.Context, email string) (User, error) {
	return f.user, f.err
}

func (f *fakeUserStore) Create(ctx context.Context, user User) error {
	return nil
}

/********** MEMBERSHIP STORE **********/
type fakeMembershipStore struct {
	membership Membership
	err        error
}

func (f *fakeMembershipStore) GetByUserID(ctx context.Context, userID string) (Membership, error) {
	if f.membership == (Membership{}) {
		return Membership{}, errs.ErrInvalidCredentials
	}
	return f.membership, nil
}

func (f *fakeMembershipStore) Create(ctx context.Context, membership Membership) error {
	return nil
}

func (f *fakeMembershipStore) GetUserFamily(ctx context.Context, familyID string) (Membership, error) {
	return f.membership, f.err
}

type fakeHasher struct{}

/********** HASHER INTERFACE **********/

func (f *fakeHasher) Compare(hash, password string) error {
	if hash == "hash" {
		return nil
	}
	return errs.ErrInvalidCredentials
}

func (f *fakeHasher) Hash(password string) (string, error) {
	return "hash", nil
}

/********** SIGNER INTERFACE **********/

type fakeSigner struct {
	token string
	err   error
}

func (f *fakeSigner) SignAccessToken(user User, m Membership) (string, error) {
	return f.token, f.err
}
