package users

const (
	UserCreatedEvent = "account.UserCreated"
)

type (
	UserCreated struct {
		ID       string
		Email    string
		Phone    string
		Name     string
		LastName string
	}
)

func (UserCreated) Key() string { return UserCreatedEvent }
