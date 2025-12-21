package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/password"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage"

	"github.com/google/uuid"
)

type FamilyStoreProvider = storage.FamilyStoreProvider
type MembershipStoreProvider = storage.MembershipStoreProvider

type RegistrationService struct {
	db   *sql.DB
	hash password.Hasher
	//functions to provide stores with attached transaction
	userStoreProvider   UserStoreProvider
	familyStoreProvider FamilyStoreProvider
	memberStoreProvider MembershipStoreProvider
}

func NewRegistrationService(db *sql.DB,
	hash password.Hasher,
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
	transaction, err := svc.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			transaction.Rollback()
		}
	}()

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

	family := domain.Family{
		ID:        uuid.NewString(),
		Name:      familyName,
		CreatedAt: time.Now(),
	}

	membership := domain.Membership{
		UserID:    user.ID,
		FamilyID:  family.ID,
		Role:      "owner",
		CreatedAt: time.Now(),
	}
	// without transaction we would skip this and pass db directly to the stores
	userStore := svc.userStoreProvider(transaction)
	familyStore := svc.familyStoreProvider(transaction)
	memberStore := svc.memberStoreProvider(transaction)

	if err = userStore.Create(ctx, user); err != nil {
		return err
	}

	if err = familyStore.Create(ctx, family); err != nil {
		return err
	}

	if err = memberStore.Create(ctx, membership); err != nil {
		return err
	}

	return transaction.Commit()
}
