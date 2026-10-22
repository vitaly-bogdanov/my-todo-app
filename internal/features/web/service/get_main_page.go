package web_service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	core_errors "github.com/vitaly-bogdanov/my-todo-app/internal/core/errors"
)

func (s *WebService) GetMainPage(ctx context.Context) ([]byte, error) {
	html, err := s.GetFile(ctx, "index.html")
	if err != nil {
		return nil, fmt.Errorf("failed to get main page: %w", err)
	}

	return html, nil
}

func (s *WebService) GetFile(ctx context.Context, name string) ([]byte, error) {
	root := os.Getenv("FRONTEND_DIST_PATH")
	if root == "" {
		return nil, fmt.Errorf("FRONTEND_DIST_PATH is empty: %w", core_errors.InternalServerError)
	}

	rootPath, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve frontend path: %w", err)
	}

	relativePath := strings.TrimPrefix(filepath.Clean("/"+name), "/")
	filePath := filepath.Join(rootPath, relativePath)
	if filePath != rootPath && !strings.HasPrefix(filePath, rootPath+string(os.PathSeparator)) {
		return nil, fmt.Errorf("invalid frontend file path: %w", core_errors.BadRequest)
	}

	file, err := s.webRepository.GetFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get frontend file: %w", err)
	}

	return file, nil
}
