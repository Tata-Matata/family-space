package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/password"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type User = domain.User
type UserStore = storage.UserStore
type UserStoreProvider = storage.UserStoreProvider

type LoginService struct {
	userStoreProvider  UserStoreProvider
	membershipProvider MembershipStoreProvider
	hash               password.Hasher
	db                 *sql.DB
}

func NewLoginService(db *sql.DB,
	hash password.Hasher,
	userStoreProvider UserStoreProvider,
	memberStoreProvider MembershipStoreProvider,
) *LoginService {
	return &LoginService{
		db:                 db,
		hash:               hash,
		userStoreProvider:  userStoreProvider,
		membershipProvider: memberStoreProvider,
	}
}

func (svc *LoginService) Login(
	ctx context.Context,
	email string,
	password string,
) (User, error) {

	// here the decision is made to use a transaction
	tx, err := svc.db.BeginTx(ctx, nil)
	if err != nil {
		return User{}, err
	}
	defer tx.Rollback()

	userStore := svc.userStoreProvider(tx)
	user, err := userStore.GetByEmail(ctx, email)
	if err != nil {
		// user not found OR DB error
		return User{}, ErrInvalidCredentials
	}

	if err := svc.hash.Compare(user.PasswordHash, password); err != nil {
		return User{}, ErrInvalidCredentials
	}

	membershipStore := svc.membershipProvider(tx)
	_, err = membershipStore.GetByUserID(ctx, user.ID)
	if err != nil {
		return User{}, ErrInvalidCredentials
	}

	//end of transaction
	tx.Commit()

	return User{}, nil
}
