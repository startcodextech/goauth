package application

import (
	"context"
	"github.com/startcodextech/goauth/users/internal/application/commands"
	"github.com/startcodextech/goauth/users/internal/domain/users"
	"github.com/startcodextech/goevents/ddd"
)

type (
	Commands interface {
		CreateUser(ctx context.Context, cmd commands.CreateUser) error
	}

	App interface {
		Commands
	}

	appCommands struct {
		commands.CreateUserHandler
	}

	Application struct {
		appCommands
	}
)

var _ App = (*Application)(nil)

func New(users users.UserRepository, publisher ddd.EventPublisher[ddd.Event]) *Application {
	return &Application{
		appCommands: appCommands{
			CreateUserHandler: commands.NewCreateUserHandler(users, publisher),
		},
	}
}
