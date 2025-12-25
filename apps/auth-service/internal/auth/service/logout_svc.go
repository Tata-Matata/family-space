package service

import (
	"context"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/refresh"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage"
)

// Logout service only revokes refresh tokens; access tokens expire automatically
type LogoutService struct {
	transactionMgr     TransactionMgr
	refreshTokenStore  storage.RefreshTokenStoreProvider
	refreshTokenHasher refresh.RefreshTokenHasher
}

type TransactionMgr = storage.TransactionMgr

func NewLogoutService(
	transactionMgr TransactionMgr,
	refreshStore storage.RefreshTokenStoreProvider,
	hasher refresh.RefreshTokenHasher,
) *LogoutService {
	return &LogoutService{
		transactionMgr:     transactionMgr,
		refreshTokenStore:  refreshStore,
		refreshTokenHasher: hasher,
	}
}

func (svc *LogoutService) Logout(
	ctx context.Context,
	rawToken string,
) (err error) {

	exec, finish, err := svc.transactionMgr.BeginTransaction(ctx, false)
	if err != nil {
		return err
	}
	defer func() {
		finish(err)
	}()

	hash, err := svc.refreshTokenHasher.Hash(rawToken)
	if err != nil {
		// invalid token → already effectively logged out
		return nil
	}

	store := svc.refreshTokenStore(exec)

	stored, err := store.GetByHash(ctx, hash)
	// idempotent logout. If we can't find the token, consider it already revoked
	if err != nil {
		return nil
	}

	if err := store.Revoke(ctx, stored.ID); err != nil {
		// DB write failure → must surface
		return err
	}

	return nil
}
