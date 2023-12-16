package account

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	account2 "github.com/startcodextech/goauth/internal/application/cqrs/events/account"
	"github.com/startcodextech/goauth/internal/domain/account"
	"github.com/startcodextech/goauth/proto"
)

// CreateUserHandler is a command to create a new user
type CreateUserHandler struct {
	eventBus *cqrs.EventBus
	service  account.UserService
	logger   watermill.LoggerAdapter
}

// ensure CreateUserHandler implements cqrs.Command interface
var _ cqrs.CommandHandler = (*CreateUserHandler)(nil)

// HandlerName returns the command's name
func (c CreateUserHandler) HandlerName() string {
	return "CreateUserHandler"
}

// NewCommand returns a new CreateUserHandler
func (c CreateUserHandler) NewCommand() interface{} {
	return &proto.CreateUser{}
}

// Handle handles the command
func (c CreateUserHandler) Handle(ctx context.Context, command interface{}) error {
	cmd := command.(*proto.CreateUser)

	id, err := c.service.Create(ctx, account.UserRegisterDto{
		Name:        cmd.Name,
		Lastname:    cmd.LastName,
		Email:       cmd.Email,
		Password:    cmd.Password,
		FacebookID:  cmd.FacebookId,
		GoogleID:    cmd.GoogleId,
		AppleID:     cmd.AppleId,
		MicrosoftID: cmd.MicrosoftId,
	})
	if err != nil {
		err := c.eventBus.Publish(ctx, &account2.UserCreatedFailed{
			Email: cmd.Email,
			Error: err.Error(),
		})
		if err != nil {
			c.logger.Error("Failed to publish event", err, nil)
			return err
		}

		c.logger.Error("Failed to create user", err, nil)
		return nil
	}

	c.logger.Info("User created", watermill.LogFields{
		"user_id": id,
	})

	err = c.eventBus.Publish(ctx, &account2.UserCreated{
		Id:       id,
		Name:     cmd.Name,
		LastName: cmd.LastName,
		Email:    cmd.Email,
	})
	if err != nil {
		c.logger.Error("Failed to publish event", err, nil)
		return err
	}

	return nil
}
