package service

import (
	"context"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/jwt"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/password"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
	errs "github.com/Tata-Matata/family-space/apps/auth-service/internal/errors"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage"
)

type User = domain.User
type Membership = domain.Membership
type UserStore = storage.UserStore
type UserStoreProvider = storage.UserStoreProvider
type TransactionManager = storage.TransactionMgr

type LoginService struct {
	userStoreProvider  UserStoreProvider
	membershipProvider MembershipStoreProvider
	hash               password.PasswordHasher
	db                 TransactionManager
	tokenSigner        jwt.TokenSigner
}

func NewLoginService(
	db TransactionManager,
	hash password.PasswordHasher,
	userStoreProvider UserStoreProvider,
	memberStoreProvider MembershipStoreProvider,
	tokenSigner jwt.TokenSigner,
) *LoginService {
	return &LoginService{
		db:                 db,
		hash:               hash,
		userStoreProvider:  userStoreProvider,
		membershipProvider: memberStoreProvider,
		tokenSigner:        tokenSigner,
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

	// here the decision is made to use a transaction,
	// i.e. exec is *sql.Tx (transactional)
	// non-transactional path would be exec := svc.db
	// (svc.db implements SQLExecutor)
	exec, finish, err := svc.db.BeginTransaction(ctx, true)
	if err != nil {
		return User{}, Membership{}, err
	}
	// commit or rollback at the end, depending on error presence
	defer func() {
		finish(err)
	}()

	// USER retrieval
	userStore := svc.userStoreProvider(exec)
	user, err := userStore.GetByEmail(ctx, email)
	if err != nil {
		// hide “user not found” vs “wrong password” distinction
		return User{}, Membership{}, errs.ErrInvalidCredentials
	}

	if err := svc.hash.Compare(user.PasswordHash, password); err != nil {
		return User{}, Membership{}, errs.ErrInvalidCredentials
	}

	// MEMBERSHIP retrieval
	membershipStore := svc.membershipProvider(exec)
	membership, err := membershipStore.GetByUserID(ctx, user.ID)
	if err != nil {
		return User{}, Membership{}, errs.ErrInvalidCredentials
	}

	return user, membership, nil
}
