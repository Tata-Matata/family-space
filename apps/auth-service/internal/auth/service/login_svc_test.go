package service_test

import (
	"context"
	"testing"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/service"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage"
)

func TestLoginService_Success(test *testing.T) {
	loginSvc := service.NewLoginService(
		&fakeDB{
			exec: &fakeSQLExecutor{},
		}, // db unused in unit test
		&fakeHasher{},
		func(exec storage.SQLExecutor) storage.UserStore {
			return &fakeUserStore{
				user: User{
					ID:           "u1",
					PasswordHash: "hash",
					Email:        "a@b.com",
				},
			}
		},
		func(exec storage.SQLExecutor) storage.MembershipStore {
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
