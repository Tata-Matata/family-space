package storage

import (
	"context"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
)

// can be implemented by both sql.DB and sql.Tx (inside a transaction)
// just capabilities needed by the stores, no implementation details
type FamilyStore interface {
	Create(ctx context.Context, family domain.Family) error
}
