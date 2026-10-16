package users_repository_postgres

import core_domain "github.com/vitaly-bogdanov/my-todo-app/internal/core/domain"

type UserModel struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

func userDomainFromModel(userModel UserModel) core_domain.User {
	return core_domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
	)
}

func userDomainsFromModels(userModels []UserModel) []core_domain.User {
	userDomains := make([]core_domain.User, len(userModels))

	for i, userModel := range userModels {
		userDomains[i] = userDomainFromModel(userModel)
	}
	return userDomains
}
