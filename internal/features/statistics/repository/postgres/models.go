package statistics_repository_postgres

import (
	"time"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
)

type TaskModel struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserID int
}

func taskDomainFromModel(taskModel TaskModel) core_domain.Task {
	return core_domain.NewTask(
		taskModel.ID,
		taskModel.Version,
		taskModel.Title,
		taskModel.Description,
		taskModel.Completed,
		taskModel.CreatedAt,
		taskModel.CompletedAt,
		taskModel.AuthorUserID,
	)
}

func taskDomainsFromModels(taskModels []TaskModel) []core_domain.Task {
	taskDomains := make([]core_domain.Task, len(taskModels))

	for i, taskModel := range taskModels {
		taskDomains[i] = taskDomainFromModel(taskModel)
	}
	return taskDomains
}
