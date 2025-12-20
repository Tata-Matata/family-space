package storage

import (
	"context"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
)

type User = domain.User

type UserStore interface {
	CreateUser(ctx context.Context, user User) error
	GetUserByEmail(ctx context.Context, email string) (User, error)
}
