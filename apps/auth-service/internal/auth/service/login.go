package service

import (
	"context"
	"errors"

	"github.com/Tata-Matata/family-space/apps/auth-service/internal/auth/password"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/domain"
	"github.com/Tata-Matata/family-space/apps/auth-service/internal/storage"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type User = domain.User
type UserStore = storage.UserStore

type LoginService struct {
	users  UserStore
	hasher password.Hasher
}

func (svc *LoginService) Login(
	ctx context.Context,
	email string,
	password string,
) (User, error) {

	user, err := svc.users.GetUserByEmail(ctx, email)
	if err != nil {
		// user not found OR DB error
		return User{}, ErrInvalidCredentials
	}

	if err := svc.hasher.Compare(user.PasswordHash, password); err != nil {
		return User{}, ErrInvalidCredentials
	}

	return user, nil
}
