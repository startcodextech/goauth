package account

import "github.com/modernice/goes/codec"

const (
	EventUserCreated = "account.user.created"
)

var UserEvents = [...]string{
	EventUserCreated,
}

type (
	UserCreated struct {
		UserID       string
		Name         string
		Lastname     string
		PasswordHash string
		Email        string
		PhotoURL     string
	}
)

func UserRegisterEvents(r codec.Registerer) {
	codec.Register[UserCreated](r, EventUserCreated)
}
