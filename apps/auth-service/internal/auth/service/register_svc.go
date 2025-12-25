package service

import (
	"context"
	"time"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/password"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage"

	"github.com/google/uuid"
)

type FamilyStoreProvider = storage.FamilyStoreProvider
type MembershipStoreProvider = storage.MembershipStoreProvider

type RegistrationService struct {
	db   TransactionManager
	hash password.PasswordHasher
	//functions to provide stores with attached transaction
	userStoreProvider   UserStoreProvider
	familyStoreProvider FamilyStoreProvider
	memberStoreProvider MembershipStoreProvider
}

func NewRegistrationService(
	db TransactionManager,
	hash password.PasswordHasher,
	userStore UserStoreProvider,
	familyStore FamilyStoreProvider,
	memberStore MembershipStoreProvider) *RegistrationService {
	return &RegistrationService{
		db:                  db,
		hash:                hash,
		userStoreProvider:   userStore,
		familyStoreProvider: familyStore,
		memberStoreProvider: memberStore,
	}
}

func (svc *RegistrationService) Register(
	ctx context.Context,
	email string,
	password string,
	familyName string,
) error {
	// start a transaction; here the decision is made to use a transaction
	exec, finish, err := svc.db.BeginTransaction(ctx, false)
	if err != nil {
		return err
	}
	// commit or rollback at the end, depending on error presence
	defer func() {
		finish(err)
	}()

	//USER creation with hashed password
	hash, err := svc.hash.Hash(password)
	if err != nil {
		return err
	}

	user := domain.User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
	}

	// without transaction we would skip this and pass db directly to the stores
	userStore := svc.userStoreProvider(exec)

	// create all entities within the transaction
	if err = userStore.Create(ctx, user); err != nil {
		return err
	}

	//FAMILY
	family := domain.Family{
		ID:        uuid.NewString(),
		Name:      familyName,
		CreatedAt: time.Now(),
	}
	familyStore := svc.familyStoreProvider(exec)
	if err = familyStore.Create(ctx, family); err != nil {
		return err
	}

	// MEMBERSHIP
	membership := domain.Membership{
		UserID:    user.ID,
		FamilyID:  family.ID,
		Role:      "owner",
		CreatedAt: time.Now(),
	}

	memberStore := svc.memberStoreProvider(exec)
	if err = memberStore.Create(ctx, membership); err != nil {
		return err
	}

	return nil
}
