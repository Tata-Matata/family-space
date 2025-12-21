package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/jwt"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/password"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type User = domain.User
type Membership = domain.Membership
type UserStore = storage.UserStore
type UserStoreProvider = storage.UserStoreProvider

type LoginService struct {
	userStoreProvider  UserStoreProvider
	membershipProvider MembershipStoreProvider
	hash               password.Hasher
	db                 *sql.DB
	tokenSigner        jwt.TokenSigner
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
) (string, error) {
	user, membership, err := svc.authenticate(ctx, email, password)
	if err != nil {
		return "", err
	}

	token, err := svc.tokenSigner.SignAccessToken(user, membership)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (svc *LoginService) authenticate(ctx context.Context, email string, password string) (User, Membership, error) {

	// here the decision is made to use a transaction
	tx, err := svc.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return User{}, Membership{}, err
	}
	defer tx.Rollback()

	userStore := svc.userStoreProvider(tx)
	user, err := userStore.GetByEmail(ctx, email)
	if err != nil {
		// hide “user not found” vs “wrong password” distinction
		return User{}, Membership{}, ErrInvalidCredentials
	}

	if err := svc.hash.Compare(user.PasswordHash, password); err != nil {
		return User{}, Membership{}, ErrInvalidCredentials
	}

	membershipStore := svc.membershipProvider(tx)
	membership, err := membershipStore.GetByUserID(ctx, user.ID)
	if err != nil {
		return User{}, Membership{}, ErrInvalidCredentials
	}

	if err := tx.Commit(); err != nil {
		return User{}, Membership{}, err
	}

	return user, membership, nil
}
