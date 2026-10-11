package users_service

import (
	"context"
	"fmt"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
)

func (s *UserService) GetUser(ctx context.Context, userID int) (core_domain.User, error) {
	user, err := s.userRepository.GetUser(ctx, userID)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
