package service_test

import (
	"context"
	"testing"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/service"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage"
)

func TestRegistrationService_Success(test *testing.T) {
	userStore := &fakeUserStore{}
	memberStore := &fakeMembershipStore{}
	familyStore := &fakeFamilyStore{}
	hasher := &fakeHasher{}

	regSvc := service.NewRegistrationService(
		&fakeDB{
			exec: &fakeSQLExecutor{},
		}, // db unused in unit test
		hasher,
		func(exec storage.SQLExecutor) UserStore {
			return userStore
		},
		func(exec storage.SQLExecutor) FamilyStore {
			return familyStore
		},
		func(exec storage.SQLExecutor) MembershipStore {
			return memberStore
		},
	)

	err := regSvc.Register(context.Background(), "a@b.com", "hash", "FamilyName")
	if err != nil {
		test.Fatalf("unexpected error")
	}

	if !hasher.called {
		test.Fatalf("password was not hashed")
	}

	if !userStore.called {
		test.Fatalf("user was not created")
	}

	if !familyStore.called {
		test.Fatalf("family was not created")
	}

	if !memberStore.called {
		test.Fatalf("membership was not created")
	}
}
