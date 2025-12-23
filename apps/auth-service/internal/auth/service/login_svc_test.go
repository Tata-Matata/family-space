package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/service"
	errs "github.com/Tata-Matata/family-space/apps/auth-service/internal/errors"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage"
)

func TestLoginService_Success(test *testing.T) {

	loginSvc := service.NewLoginService(
		&fakeDB{
			exec: &fakeSQLExecutor{},
		}, // db unused in unit test
		&fakeHasher{hash: HASH},
		func(exec storage.SQLExecutor) UserStore {
			return &fakeUserStore{
				user: User{
					ID:           "u1",
					PasswordHash: HASH,
					Email:        "a@b.com",
				},
			}
		},
		func(exec storage.SQLExecutor) MembershipStore {
			return &fakeMembershipStore{
				membership: Membership{
					UserID:   "u1",
					FamilyID: "f1",
					Role:     "admin",
				},
			}
		}, &fakeSigner{token: JWTToken},
	)

	token, err := loginSvc.Login(context.Background(), "a@b.com", "pw")
	if err != nil {
		test.Fatalf("unexpected Login error: %v", err)
	}

	if token != JWTToken {
		test.Fatalf("expected token %s, got '%s'", JWTToken, token)
	}
}

func TestLoginService_InvalidPassword(test *testing.T) {
	loginSvc := service.NewLoginService(
		&fakeDB{
			exec: &fakeSQLExecutor{},
		}, // db unused in unit test
		&fakeHasher{hash: HASH},
		func(exec storage.SQLExecutor) UserStore {
			return &fakeUserStore{
				user: User{PasswordHash: "wronghash"},
			}
		},
		func(exec storage.SQLExecutor) MembershipStore {
			return &fakeMembershipStore{}
		},

		&fakeSigner{},
	)

	_, err := loginSvc.Login(context.Background(), "a@b.com", "pw")
	if !errors.Is(err, errs.ErrInvalidCredentials) {
		test.Fatalf("expected %v, but got: %v", errs.ErrInvalidCredentials, err)
	}
}

func TestLoginService_UserNotFound(test *testing.T) {
	loginSvc := service.NewLoginService(
		&fakeDB{
			exec: &fakeSQLExecutor{},
		}, // db unused in unit test
		&fakeHasher{hash: HASH},
		func(exec storage.SQLExecutor) UserStore {
			//whatever error userStore returns, login svc should obscure it as invalid credentials
			return &fakeUserStore{
				err: errs.ErrNotFound,
			}
		},
		func(exec storage.SQLExecutor) MembershipStore {
			return &fakeMembershipStore{}
		},
		&fakeSigner{token: JWTToken},
	)

	_, err := loginSvc.Login(context.Background(), "a@b.com", "pw")
	if !errors.Is(err, errs.ErrInvalidCredentials) {
		test.Fatalf("expected %v, but got: %v", errs.ErrInvalidCredentials, err)
	}
}

func TestLoginService_MembershipNotFound(test *testing.T) {
	loginSvc := service.NewLoginService(
		&fakeDB{
			exec: &fakeSQLExecutor{},
		}, // db unused in unit test
		&fakeHasher{hash: HASH},
		func(exec storage.SQLExecutor) UserStore {
			return &fakeUserStore{user: User{
				ID:           "u1",
				PasswordHash: HASH,
				Email:        "a@b.com",
			}}
		},
		func(exec storage.SQLExecutor) MembershipStore {
			//whatever error the store returns, login svc should obscure it as invalid credentials
			return &fakeMembershipStore{
				err: errs.ErrNotFound,
			}
		},
		&fakeSigner{token: JWTToken},
	)

	_, err := loginSvc.Login(context.Background(), "a@b.com", "pw")
	if !errors.Is(err, errs.ErrInvalidCredentials) {
		test.Fatalf("expected %v, but got: %v", errs.ErrInvalidCredentials, err)
	}
}
