package users_repository_postgres

import core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"

type UserModel struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

func userDomainsFromModels(userModels []UserModel) []core_domain.User {
	userDomains := make([]core_domain.User, len(userModels))

	for i, userModel := range userModels {
		userDomains[i] = core_domain.NewUser(
			userModel.ID,
			userModel.Version,
			userModel.FullName,
			userModel.PhoneNumber,
		)
	}
	return userDomains
}
