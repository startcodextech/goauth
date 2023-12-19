package account

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/startcodextech/goauth/internal/domain/account"
	"github.com/startcodextech/goauth/proto"
	"go.uber.org/zap"
)

// CreateUserHandler is a command to create a new user
type CreateUserHandler struct {
	eventBus *cqrs.EventBus
	service  account.UserService
	logger   *zap.Logger
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
		err := c.eventBus.Publish(ctx, &proto.EventUserCreatedFailed{
			Email: cmd.Email,
			Error: err.Error(),
		})
		if err != nil {
			c.logger.Error("Failed to publish event.go", zap.Error(err))
			return err
		}

		c.logger.Error(
			"Failed to create user",
			zap.String("email", cmd.Email),
			zap.Error(err),
		)
		return nil
	}

	c.logger.Info(
		"User created",
		zap.String("user_id", id),
	)

	err = c.eventBus.Publish(ctx, &proto.EventUserCreated{
		Id:       id,
		Name:     cmd.Name,
		LastName: cmd.LastName,
		Email:    cmd.Email,
	})
	if err != nil {
		c.logger.Error("Failed to publish event.go", zap.Error(err))
		return err
	}

	return nil
}
