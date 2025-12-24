package storage

import (
	"context"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/refresh"
)

type RefreshToken = refresh.RefreshToken

type RefreshTokenStore interface {
	Create(ctx context.Context, token RefreshToken) error
	GetByHash(ctx context.Context, hash string) (RefreshToken, error)
	Revoke(ctx context.Context, id string) error
}
