package service

import (
	"context"
	"time"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/jwt"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/refresh"
	errs "github.com/Tata-Matata/family-space/apps/auth-service/internal/errors"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage"
	"github.com/google/uuid"
)

// RefreshService handles refresh token rotation and access token issuance
type RefreshService struct {
	transactionMgr     storage.TransactionMgr
	refreshTokenStore  storage.RefreshTokenStoreProvider
	userStoreProvider  storage.UserStoreProvider
	membershipProvider storage.MembershipStoreProvider
	refreshTokenHasher refresh.RefreshTokenHasher
	refreshTokenGen    refresh.RefreshTokenGenerator
	tokenSigner        jwt.TokenSigner
	accessTokenTTL     time.Duration
}

func NewRefreshService(
	transactionMgr storage.TransactionMgr,
	refreshTokenStore storage.RefreshTokenStoreProvider,
	userStore storage.UserStoreProvider,
	membershipStore storage.MembershipStoreProvider,
	refreshHasher refresh.RefreshTokenHasher,
	refreshGen refresh.RefreshTokenGenerator,
	signer jwt.TokenSigner,
	accessTTL time.Duration,
) *RefreshService {
	return &RefreshService{
		transactionMgr:     transactionMgr,
		refreshTokenStore:  refreshTokenStore,
		userStoreProvider:  userStore,
		membershipProvider: membershipStore,
		refreshTokenHasher: refreshHasher,
		refreshTokenGen:    refreshGen,
		tokenSigner:        signer,
		accessTokenTTL:     accessTTL,
	}
}

func (svc *RefreshService) Refresh(
	ctx context.Context,
	rawRefreshToken string,
) (newAccessToken string, newRefreshToken string, err error) {

	exec, finish, err := svc.transactionMgr.BeginTransaction(ctx, false)
	if err != nil {
		return "", "", err
	}
	defer func() {
		err = finish(err)
	}()

	refreshStore := svc.refreshTokenStore(exec)

	// 1. Hash incoming refresh token
	hash, err := svc.refreshTokenHasher.Hash(rawRefreshToken)
	if err != nil {
		return "", "", errs.ErrInvalidRefreshToken
	}

	// 2. Load refresh token record
	stored, err := refreshStore.GetByHash(ctx, hash)
	if err != nil {
		return "", "", errs.ErrInvalidRefreshToken
	}

	// 3. Validate refresh token
	if stored.RevokedAt != nil {
		return "", "", errs.ErrInvalidRefreshToken
	}

	if time.Now().After(stored.ExpiresAt) {
		return "", "", errs.ErrInvalidRefreshToken
	}

	// 4. Revoke old refresh token (rotation)
	if err := refreshStore.Revoke(ctx, stored.ID); err != nil {
		return "", "", err
	}

	// 5. Load user + membership
	userStore := svc.userStoreProvider(exec)
	user, err := userStore.GetById(ctx, stored.UserID)
	if err != nil {
		return "", "", err
	}

	membershipStore := svc.membershipProvider(exec)
	membership, err := membershipStore.GetByUserID(ctx, user.ID)
	if err != nil {
		return "", "", err
	}

	// 6. Issue new access token
	accessToken, err := svc.tokenSigner.SignAccessToken(user, membership)
	if err != nil {
		return "", "", err
	}

	// 7. Generate + store new refresh token
	refreshToken, err := svc.refreshTokenGen.Generate()
	if err != nil {
		return "", "", err
	}

	// hash new refresh token
	newHash, err := svc.refreshTokenHasher.Hash(refreshToken)
	if err != nil {
		return "", "", err
	}

	// store new refresh token in DB
	tokenStoredInDb := refresh.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		TokenHash: newHash,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	if err := refreshStore.Create(ctx, tokenStoredInDb); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
