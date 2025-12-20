package storage

import (
	"context"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
)

type Membership = domain.Membership

type MembershipStore interface {
	AddMembership(ctx context.Context, m Membership) error
	GetUserFamily(ctx context.Context, userID string) (Membership, error)
}
