package users_transport_http

import core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"

type UserDTO struct {
	ID          int     `json:"id"           example:"1"`
	Version     int     `json:"version"      example:"1"`
	FullName    string  `json:"full_name"    example:"John Doe"`
	PhoneNumber *string `json:"phone_number" example:"+1234567890"`
}

func userDtoFromDomain(domain core_domain.User) UserDTO {
	return UserDTO{
		ID:          domain.ID,
		Version:     domain.Version,
		FullName:    domain.FullName,
		PhoneNumber: domain.PhoneNumber,
	}
}

func userDtosFromDomains(users []core_domain.User) []UserDTO {
	usersDTO := make([]UserDTO, len(users))

	for i, user := range users {
		usersDTO[i] = userDtoFromDomain(user)
	}

	return usersDTO
}
