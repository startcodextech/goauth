package pb

import "github.com/startcodextech/goevents/registry"

const (
	UserAggregateChannel = "goauth.users.events.User"

	UserCreatedEvent = "usersapi.UserCreated"

	CommandChannel = "goauth.users.commands"
)

func Registrations(reg registry.Registry) (err error) {
	return nil
}
