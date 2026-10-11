package users_service

import (
	"context"
	"fmt"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
	"golang.org/x/sync/errgroup"
)

func (s *UserService) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]core_domain.User, int, error) {
	if limit != nil && *limit <= 0 {
		return nil, 0, fmt.Errorf("limit must be greater than 0")
	}

	if offset != nil && *offset < 0 {
		return nil, 0, fmt.Errorf("offset must be greater than or equal to 0")
	}

	g, gctx := errgroup.WithContext(ctx)

	var users []core_domain.User
	g.Go(func() error {
		var err error
		users, err = s.userRepository.GetUsers(gctx, limit, offset)
		if err != nil {
			return fmt.Errorf("failed to get users: %w", err)
		}
		return nil
	})

	var total int
	g.Go(func() error {
		var err error
		total, err = s.userRepository.CountUsers(gctx)
		if err != nil {
			return fmt.Errorf("failed to count users: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
