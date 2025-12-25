package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/refresh"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/service"
	errs "github.com/Tata-Matata/family-space/apps/auth-service/internal/errors"
)

func TestLogoutService_Success(test *testing.T) {
	store := &fakeRefreshTokenStore{
		token: refresh.RefreshToken{ID: "token-id"},
	}

	svc := service.NewLogoutService(
		&fakeDB{},
		refreshStoreProvider(store),
		&fakeRefreshTokenHasher{hash: "hash"},
	)

	err := svc.Logout(context.Background(), "raw-token")
	if err != nil {
		test.Fatalf("unexpected error: %v", err)
	}

	if !store.revokeCalled {
		test.Fatalf("expected revoke to be called")
	}
}

func TestLogoutService_TokenNotFound(test *testing.T) {
	store := &fakeRefreshTokenStore{
		getErr: errs.ErrNotFound,
	}

	svc := service.NewLogoutService(
		&fakeDB{},
		refreshStoreProvider(store),
		&fakeRefreshTokenHasher{hash: "hash"},
	)

	err := svc.Logout(context.Background(), "raw-token")
	if err != nil {
		test.Fatalf("expected no error. Got %v", err)
	}
}

func TestLogoutService_InvalidToken(test *testing.T) {
	hasher := &fakeRefreshTokenHasher{
		err: errors.New("invalid"),
	}

	svc := service.NewLogoutService(
		&fakeDB{},
		refreshStoreProvider(&fakeRefreshTokenStore{}),
		hasher,
	)

	err := svc.Logout(context.Background(), "bad-token")
	if err != nil {
		test.Fatalf("expected no error. Got %v", err)
	}
}

func TestLogoutService_RevokeFailure(test *testing.T) {
	store := &fakeRefreshTokenStore{
		token:     refresh.RefreshToken{ID: "token-id"},
		revokeErr: errors.New("db failure"),
	}

	svc := service.NewLogoutService(
		&fakeDB{},
		refreshStoreProvider(store),
		&fakeRefreshTokenHasher{hash: "hash"},
	)

	err := svc.Logout(context.Background(), "raw-token")
	if err == nil {
		test.Fatalf("expected error")
	}
}

func TestLogoutService_BeginTransactionFailure(test *testing.T) {
	svc := service.NewLogoutService(
		&fakeDB{beginErr: errors.New("transaction failed at the start")},
		refreshStoreProvider(&fakeRefreshTokenStore{}),
		&fakeRefreshTokenHasher{},
	)

	err := svc.Logout(context.Background(), "raw-token")
	if err == nil {
		test.Fatalf("expected error")
	}
}
