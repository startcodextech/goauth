package account

import (
	"context"
	"github.com/google/uuid"
	"github.com/modernice/goes/aggregate"
	"github.com/modernice/goes/codec"
	"github.com/modernice/goes/command"
	"github.com/modernice/goes/command/handler"
)

const (
	UserCreatedCmd = "account.user.create"
)

func CreateUser(userID uuid.UUID, user UserCreateDto) command.Cmd[UserCreateDto] {
	return command.New[UserCreateDto](UserCreatedCmd, user, command.Aggregate(UserAggregate, userID))
}

func userRegisterCommands(r codec.Registerer) {
	codec.Register[UserCreateDto](r, UserCreatedCmd)
}

func UserHandleCommands(ctx context.Context, bus command.Bus, repo aggregate.Repository) <-chan error {
	return handler.New(UserNew, repo, bus).MustHandle(ctx)
}
