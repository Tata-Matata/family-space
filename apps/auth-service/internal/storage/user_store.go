package storage

import (
	"context"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
)

// can be implemented by both sql.DB and sql.Tx (inside a transaction)
// just capabilities needed by the stores, no implementation details
type UserStore interface {
	Create(ctx context.Context, user domain.User) error
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetById(ctx context.Context, id string) (domain.User, error)
}
