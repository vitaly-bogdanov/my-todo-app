package web_repository_file_system

import (
	"errors"
	"fmt"
	"os"

	core_errors "github.com/vitaly-bogdanov/my-todo-app/internal/core/errors"
)

func (r *WebRepository) GetFile(filePath string) ([]byte, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("file not found: %w: %s", core_errors.NotFound, filePath)
		}

		return nil, fmt.Errorf("failed to read file: %w: %s", err, filePath)
	}

	return file, nil
}
