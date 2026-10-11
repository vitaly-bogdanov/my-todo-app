package users_service

import (
	"context"
	"fmt"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
)

func (s *UserService) PatchUser(
	ctx context.Context,
	userID int,
	userPatch core_domain.UserPatch,
) (core_domain.User, error) {
	user, err := s.userRepository.GetUser(ctx, userID)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	if err := user.ApplyPatch(userPatch); err != nil {
		return core_domain.User{}, fmt.Errorf("failed to apply patch: %w", err)
	}

	patchedUser, err := s.userRepository.PatchUser(ctx, userID, user)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("failed to patch user: %w", err)
	}

	return patchedUser, nil
}
