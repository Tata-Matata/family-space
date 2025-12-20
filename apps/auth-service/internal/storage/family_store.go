package storage

import (
	"context"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
)

type Family = domain.Family

type FamilyStore interface {
	CreateFamily(ctx context.Context, family Family) error
}
