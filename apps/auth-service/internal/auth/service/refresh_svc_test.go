package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/refresh"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/service"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
	errs "github.com/Tata-Matata/family-space/apps/auth-service/internal/errors"
)

func TestRefreshService_Success(test *testing.T) {
	now := time.Now()

	refreshStore := &fakeRefreshTokenStore{
		token: refresh.RefreshToken{
			ID:        "old-id",
			UserID:    "user-1",
			ExpiresAt: now.Add(time.Hour),
		},
	}

	svc := service.NewRefreshService(
		&fakeDB{},
		refreshStoreProvider(refreshStore),
		userStoreProvider(&fakeUserStore{
			user: domain.User{ID: "user-1"},
		}),
		membershipStoreProvider(&fakeMembershipStore{
			membership: domain.Membership{UserID: "user-1"},
		}),
		&fakeRefreshTokenHasher{hash: "hash"},
		&fakeRefreshTokenGenerator{token: "new-refresh"},
		&fakeSigner{token: "new-access"},
		15*time.Minute,
	)

	access, refresh, err := svc.Refresh(context.Background(), "raw-token")
	if err != nil {
		test.Fatalf("unexpected error: %v", err)
	}

	if access != "new-access" {
		test.Fatalf("unexpected access token")
	}
	if refresh != "new-refresh" {
		test.Fatalf("unexpected refresh token")
	}
	if !refreshStore.revokeCalled {
		test.Fatalf("expected revoke to be called")
	}
	if !refreshStore.createCalled {
		test.Fatalf("expected create to be called")
	}
}

func TestRefreshService_ExpiredToken(test *testing.T) {
	refreshStore := &fakeRefreshTokenStore{
		token: refresh.RefreshToken{
			ID:        "id",
			UserID:    "user",
			ExpiresAt: time.Now().Add(-time.Minute),
		},
	}

	svc := service.NewRefreshService(
		&fakeDB{},
		refreshStoreProvider(refreshStore),
		userStoreProvider(&fakeUserStore{}),
		membershipStoreProvider(&fakeMembershipStore{}),
		&fakeRefreshTokenHasher{hash: "hash"},
		&fakeRefreshTokenGenerator{},
		&fakeSigner{},
		15*time.Minute,
	)

	_, _, err := svc.Refresh(context.Background(), "raw-token")
	if !errors.Is(err, errs.ErrInvalidRefreshToken) {
		test.Fatalf("expected invalid refresh token error")
	}
}

func TestRefreshService_RevokedToken(test *testing.T) {
	now := time.Now()

	refreshStore := &fakeRefreshTokenStore{
		token: refresh.RefreshToken{
			ID:        "id",
			UserID:    "user",
			ExpiresAt: now.Add(time.Hour),
			RevokedAt: &now,
		},
	}

	svc := service.NewRefreshService(
		&fakeDB{},
		refreshStoreProvider(refreshStore),
		userStoreProvider(&fakeUserStore{}),
		membershipStoreProvider(&fakeMembershipStore{}),
		&fakeRefreshTokenHasher{hash: "hash"},
		&fakeRefreshTokenGenerator{},
		&fakeSigner{},
		15*time.Minute,
	)

	_, _, err := svc.Refresh(context.Background(), "raw-token")
	if !errors.Is(err, errs.ErrInvalidRefreshToken) {
		test.Fatalf("expected invalid refresh token error")
	}
}

func TestRefreshService_RevokeFailure(test *testing.T) {
	refreshStore := &fakeRefreshTokenStore{
		token: refresh.RefreshToken{
			ID:        "id",
			UserID:    "user",
			ExpiresAt: time.Now().Add(time.Hour),
		},
		revokeErr: errors.New("db error"),
	}

	svc := service.NewRefreshService(
		&fakeDB{},
		refreshStoreProvider(refreshStore),
		userStoreProvider(&fakeUserStore{}),
		membershipStoreProvider(&fakeMembershipStore{}),
		&fakeRefreshTokenHasher{hash: "hash"},
		&fakeRefreshTokenGenerator{},
		&fakeSigner{},
		15*time.Minute,
	)

	_, _, err := svc.Refresh(context.Background(), "raw-token")
	if err == nil {
		test.Fatalf("expected error")
	}
}

func TestRefreshService_SignerFailure(test *testing.T) {
	refreshStore := &fakeRefreshTokenStore{
		token: refresh.RefreshToken{
			ID:        "id",
			UserID:    "user",
			ExpiresAt: time.Now().Add(time.Hour),
		},
	}

	svc := service.NewRefreshService(
		&fakeDB{},
		refreshStoreProvider(refreshStore),
		userStoreProvider(&fakeUserStore{}),
		membershipStoreProvider(&fakeMembershipStore{}),
		&fakeRefreshTokenHasher{hash: "hash"},
		&fakeRefreshTokenGenerator{token: "new"},
		&fakeSigner{err: errors.New("sign fail")},
		15*time.Minute,
	)

	_, _, err := svc.Refresh(context.Background(), "raw-token")
	if err == nil {
		test.Fatalf("expected error")
	}
}
