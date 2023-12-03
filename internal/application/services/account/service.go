package account

import "github.com/startcodextech/goauth/internal/domain/account"

// AccountServices is a struct that contains all the services
// related to the account domain
type AccountServices struct {
	user account.UserService
}

// New creates a new AccountServices struct
// with all the account services
func New(repo account.UserRepository) AccountServices {
	return AccountServices{
		user: NewUserService(repo),
	}
}

// User returns the user service, that is used to manage users
func (s AccountServices) User() account.UserService {
	return s.user
}
