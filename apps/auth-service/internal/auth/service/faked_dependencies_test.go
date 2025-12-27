package service_test

import (
	"context"
	"database/sql"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/refresh"
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

const JWTToken = "jwt.token"

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
	exec     storage.SQLExecutor
	beginErr error
}

// immitates starting db transaction but actually does nothing
func (fakeDB *fakeDB) BeginTransaction(
	ctx context.Context,
	readOnly bool,
) (storage.SQLExecutor, func(error), error) {
	if fakeDB.beginErr != nil {
		return nil, nil, fakeDB.beginErr
	}

	finish := func(err error) {
	}

	return fakeDB.exec, finish, nil
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
	return fakeUserStore.err
}
func (fakeUserStore *fakeUserStore) GetById(ctx context.Context, id string) (User, error) {
	return fakeUserStore.user, fakeUserStore.err
}

func userStoreProvider(store *fakeUserStore) storage.UserStoreProvider {
	return func(exec storage.SQLExecutor) storage.UserStore {
		return store
	}
}

/********** MEMBERSHIP STORE **********/
type fakeMembershipStore struct {
	membership Membership
	err        error
	called     bool
}

func (fakeMemStore *fakeMembershipStore) GetByUserID(ctx context.Context, userID string) (Membership, error) {

	return fakeMemStore.membership, fakeMemStore.err
}

func (fakeMemStore *fakeMembershipStore) Create(ctx context.Context, membership Membership) error {
	fakeMemStore.called = true
	return fakeMemStore.err
}

func (fakeMemStore *fakeMembershipStore) GetUserFamily(ctx context.Context, familyID string) (Membership, error) {
	return fakeMemStore.membership, fakeMemStore.err
}

func membershipStoreProvider(store *fakeMembershipStore) storage.MembershipStoreProvider {
	return func(exec storage.SQLExecutor) storage.MembershipStore {
		return store
	}
}

/*** FAKE FAMILY STORE ***/
type fakeFamilyStore struct {
	family Family
	err    error
	called bool
}

func (fakeFamilyStore *fakeFamilyStore) Create(ctx context.Context, family Family) error {
	fakeFamilyStore.called = true

	return fakeFamilyStore.err
}

/********** HASHER INTERFACE **********/
const HASH = "hash"

type fakeHasher struct {
	err    error
	hash   string
	called bool
}

func (hasher *fakeHasher) Compare(hash, password string) error {
	hasher.called = true
	if hash == HASH {
		return nil
	}
	return errs.ErrInvalidCredentials
}

func (hasher *fakeHasher) Hash(password string) (string, error) {
	hasher.called = true
	return hasher.hash, hasher.err
}

/********** SIGNER INTERFACE **********/

type fakeSigner struct {
	token string
	err   error
}

func (signer *fakeSigner) GenerateSignedAccessToken(user User, m Membership) (string, error) {
	return signer.token, signer.err
}

// ******** Logout and refresh service **********/
type fakeRefreshTokenHasher struct {
	hash string
	err  error
}

func (hasher *fakeRefreshTokenHasher) Hash(token string) (string, error) {
	if hasher.err != nil {
		return "", hasher.err
	}
	return hasher.hash, nil
}

func (hasher *fakeRefreshTokenHasher) Compare(hash, token string) error {
	return nil
}

type fakeRefreshTokenStore struct {
	token        refresh.RefreshToken
	getErr       error
	revokeErr    error
	createErr    error
	revokeCalled bool
	createCalled bool
}

func (refreshStore *fakeRefreshTokenStore) GetByHash(
	ctx context.Context,
	hash string,
) (refresh.RefreshToken, error) {
	if refreshStore.getErr != nil {
		return refresh.RefreshToken{}, refreshStore.getErr
	}
	return refreshStore.token, nil
}

func (refreshStore *fakeRefreshTokenStore) Revoke(
	ctx context.Context,
	id string,
) error {
	refreshStore.revokeCalled = true
	return refreshStore.revokeErr
}

func (refreshStore *fakeRefreshTokenStore) Create(
	ctx context.Context,
	token refresh.RefreshToken,
) error {
	refreshStore.createCalled = true
	return refreshStore.createErr
}

func refreshStoreProvider(store *fakeRefreshTokenStore) storage.RefreshTokenStoreProvider {
	return func(exec storage.SQLExecutor) storage.RefreshTokenStore {
		return store
	}
}

type fakeRefreshTokenGenerator struct {
	token string
	err   error
}

func (f *fakeRefreshTokenGenerator) Generate() (string, error) {
	if f.err != nil {
		return "", f.err
	}
	return f.token, nil
}
