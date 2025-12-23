package service_test

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
	errs "github.com/Tata-Matata/family-space/apps/auth-service/internal/errors"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage"
)

type User = domain.User
type Membership = domain.Membership
type Family = domain.Family

type UserStore = storage.UserStore
type MembershipStore = storage.MembershipStore
type FamilyStore = storage.FamilyStore

/**FAKE DATABASE FOR UNIT TESTS **/
type fakeSQLExecutor struct{}

func (fakeSqlExec *fakeSQLExecutor) ExecContext(
	ctx context.Context,
	query string,
	args ...any,
) (sql.Result, error) {
	panic("ExecContext should not be called in service unit test")
}

func (fakeSqlExec *fakeSQLExecutor) QueryRowContext(
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
func (fakeDB *fakeDB) BeginTransaction(
	ctx context.Context,
	readOnly bool,
) (storage.SQLExecutor, func(error) error, error) {
	return fakeDB.exec, func(error) error { return nil }, nil
}

/********** USER STORE **********/
type fakeUserStore struct {
	user   User
	err    error
	called bool
}

func (fakeUserStore *fakeUserStore) GetByEmail(ctx context.Context, email string) (User, error) {
	return fakeUserStore.user, fakeUserStore.err
}

func (fakeUserStore *fakeUserStore) Create(ctx context.Context, user User) error {
	fakeUserStore.called = true
	return nil
}

/********** MEMBERSHIP STORE **********/
type fakeMembershipStore struct {
	membership Membership
	err        error
	called     bool
}

func (fakeMemStore *fakeMembershipStore) GetByUserID(ctx context.Context, userID string) (Membership, error) {
	if fakeMemStore.membership == (Membership{}) {
		return Membership{}, errs.ErrInvalidCredentials
	}
	return fakeMemStore.membership, nil
}

func (fakeMemStore *fakeMembershipStore) Create(ctx context.Context, membership Membership) error {
	fakeMemStore.called = true
	return nil
}

func (fakeMemStore *fakeMembershipStore) GetUserFamily(ctx context.Context, familyID string) (Membership, error) {
	return fakeMemStore.membership, fakeMemStore.err
}

/*** FAKE FAMILY STORE ***/
type fakeFamilyStore struct {
	family Family
	err    error
	called bool
}

func (fakeFamilyStore *fakeFamilyStore) Create(ctx context.Context, family Family) error {
	fakeFamilyStore.called = true
	if fakeFamilyStore.family == (Family{}) {
		return errors.New("family not created")
	}
	return nil
}

/********** HASHER INTERFACE **********/
type fakeHasher struct {
	called bool
}

func (hasher *fakeHasher) Compare(hash, password string) error {
	hasher.called = true
	if hash == "hash" {
		return nil
	}
	return errs.ErrInvalidCredentials
}

func (hasher *fakeHasher) Hash(password string) (string, error) {
	hasher.called = true
	return "hash", nil
}

/********** SIGNER INTERFACE **********/

type fakeSigner struct {
	token string
	err   error
}

func (signer *fakeSigner) SignAccessToken(user User, m Membership) (string, error) {
	return signer.token, signer.err
}
