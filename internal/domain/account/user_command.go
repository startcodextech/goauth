package account

import (
	"github.com/google/uuid"
	"github.com/modernice/goes/command"
)

const (
	UserCreatedCmd = "account.user.create"
)

func CreateUser(userID uuid.UUID, user UserCreateDto) command.Cmd[UserCreateDto] {
	return command.New(UserCreatedCmd, user, command.Aggregate(UserAggregate, userID))
}
