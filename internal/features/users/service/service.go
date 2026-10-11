package users_service

import (
	"context"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
)

type UserService struct {
	userRepository UserRepository
}

type UserRepository interface {
	CreateUser(ctx context.Context, user core_domain.User) (core_domain.User, error)
	GetUsers(ctx context.Context, limit *int, offset *int) ([]core_domain.User, error)
	CountUsers(ctx context.Context) (int, error)
	GetUser(ctx context.Context, userID int) (core_domain.User, error)
	DeleteUser(ctx context.Context, userID int) error
	PatchUser(ctx context.Context, userID int, user core_domain.User) (core_domain.User, error)
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}
