package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/service"
	errs "github.com/Tata-Matata/family-space/apps/auth-service/internal/errors"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage"
)

type UserStore = storage.UserStore
type MembershipStore = storage.MembershipStore

func TestLoginService_Success(test *testing.T) {
	loginSvc := service.NewLoginService(
		&fakeDB{
			exec: &fakeSQLExecutor{},
		}, // db unused in unit test
		&fakeHasher{},
		func(exec storage.SQLExecutor) UserStore {
			return &fakeUserStore{
				user: User{
					ID:           "u1",
					PasswordHash: "hash",
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
		},
		&fakeSigner{token: "jwt.token"},
	)

	token, err := loginSvc.Login(context.Background(), "a@b.com", "pw")
	if err != nil {
		test.Fatalf("unexpected Login error: %v", err)
	}

	if token != "jwt.token" {
		test.Fatalf("expected token 'jwt.token', got '%s'", token)
	}
}

func TestLoginService_InvalidPassword(test *testing.T) {
	loginSvc := service.NewLoginService(
		&fakeDB{
			exec: &fakeSQLExecutor{},
		}, // db unused in unit test
		&fakeHasher{},
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
		test.Fatalf("expected ErrInvalidCredentials")
	}
}

func TestLoginService_UserNotFound(test *testing.T) {
	loginSvc := service.NewLoginService(
		&fakeDB{
			exec: &fakeSQLExecutor{},
		}, // db unused in unit test
		&fakeHasher{},
		func(exec storage.SQLExecutor) UserStore {
			return &fakeUserStore{}
		},
		func(exec storage.SQLExecutor) MembershipStore {
			return &fakeMembershipStore{}
		},
		&fakeSigner{token: "jwt.token"},
	)

	_, err := loginSvc.Login(context.Background(), "a@b.com", "pw")
	if !errors.Is(err, errs.ErrInvalidCredentials) {
		test.Fatalf("expected ErrInvalidCredentials")
	}
}

func TestLoginService_MembershipNotFound(test *testing.T) {
	loginSvc := service.NewLoginService(
		&fakeDB{
			exec: &fakeSQLExecutor{},
		}, // db unused in unit test
		&fakeHasher{},
		func(exec storage.SQLExecutor) UserStore {
			return &fakeUserStore{user: User{
				ID:           "u1",
				PasswordHash: "hash",
				Email:        "a@b.com",
			}}
		},
		func(exec storage.SQLExecutor) MembershipStore {
			return &fakeMembershipStore{}
		},
		&fakeSigner{token: "jwt.token"},
	)

	_, err := loginSvc.Login(context.Background(), "a@b.com", "pw")
	if !errors.Is(err, errs.ErrInvalidCredentials) {
		test.Fatalf("expected ErrInvalidCredentials, but got: %v", err)
	}
}
