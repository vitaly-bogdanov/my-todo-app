package tasks_transport_http

import (
	"time"

	core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"
)

type TaskDTO struct {
	ID           int        `json:"id"`
	Version      int        `json:"version"`
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
	AuthorUserID int        `json:"author_user_id"`
	CreatedAt    time.Time  `json:"created_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	Completed    bool       `json:"completed"`
}

func taskDtoFromDomain(domain core_domain.Task) TaskDTO {
	return TaskDTO{
		ID:           domain.ID,
		Version:      domain.Version,
		Title:        domain.Title,
		Description:  domain.Description,
		AuthorUserID: domain.AuthorUserID,
		CreatedAt:    domain.CreatedAt,
		CompletedAt:  domain.CompletedAt,
		Completed:    domain.Completed,
	}
}

func taskDtosFromDomains(domains []core_domain.Task) []TaskDTO {
	dto := make([]TaskDTO, len(domains))
	for i, domain := range domains {
		dto[i] = taskDtoFromDomain(domain)
	}
	return dto
}
