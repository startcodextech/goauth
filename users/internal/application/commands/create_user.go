package commands

import (
	"context"
	"github.com/pkg/errors"
	"github.com/startcodextech/goauth/users/internal/domain/users"
	"github.com/startcodextech/goevents/ddd"
)

type (
	CreateUser struct {
		ID       string
		Email    string
		Phone    string
		Password string
		Name     string
		Lastname string
	}

	CreateUserHandler struct {
		users     users.UserRepository
		publisher ddd.EventPublisher[ddd.Event]
	}
)

func NewCreateUserHandler(users users.UserRepository, publisher ddd.EventPublisher[ddd.Event]) CreateUserHandler {
	return CreateUserHandler{
		users:     users,
		publisher: publisher,
	}
}

func (h CreateUserHandler) CreateUser(ctx context.Context, cmd CreateUser) error {
	user, err := h.users.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := user.CreateUser(
		cmd.ID,
		cmd.Email,
		cmd.Phone,
		cmd.Password,
		cmd.Name,
		cmd.Lastname,
	)
	if err != nil {
		return errors.Wrap(err, "create user command")
	}

	if err = h.users.Save(ctx, user); err != nil {
		return errors.Wrap(err, "user creation")
	}

	return h.publisher.Publish(ctx, event)
}
