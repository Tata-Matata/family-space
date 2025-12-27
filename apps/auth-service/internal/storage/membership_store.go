package storage

import (
	"context"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
)

type Membership = domain.Membership

// can be implemented by both sql.DB and sql.Tx (inside a transaction)
// just capabilities needed by the stores, no implementation details
type MembershipStore interface {
	Create(ctx context.Context, m Membership) error
	GetByUserID(ctx context.Context, userID string) (Membership, error)
}
