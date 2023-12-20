package account

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/startcodextech/goauth/internal/domain/account"
	"github.com/startcodextech/goauth/proto"
	"github.com/startcodextech/goauth/util/channel"
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
	return &proto.CommandCreateUser{}
}

// Handle handles the command
func (c CreateUserHandler) Handle(ctx context.Context, command interface{}) error {
	cmd := command.(*proto.CommandCreateUser)

	c.logger.Info("creating user", zap.String("command_id", cmd.GetCommandId()))

	id, err := c.service.Create(ctx, account.UserRegisterDto{
		Name:        cmd.GetPayload().GetName(),
		Lastname:    cmd.GetPayload().GetLastName(),
		Email:       cmd.GetPayload().GetEmail(),
		Password:    cmd.GetPayload().GetPassword(),
		FacebookID:  cmd.GetPayload().GetFacebookId(),
		GoogleID:    cmd.GetPayload().GetGoogleId(),
		AppleID:     cmd.GetPayload().GetAppleId(),
		MicrosoftID: cmd.GetPayload().GetMicrosoftId(),
	})
	if err != nil {
		event := &proto.EventError{
			CommandId: cmd.GetCommandId(),
			Error:     err.Error(),
		}

		channel.Channels.Failed(cmd.GetCommandId(), event)

		err := c.eventBus.Publish(ctx, event)
		if err != nil {
			c.logger.Error("Failed to publish event.go", zap.Error(err))
			return err
		}

		c.logger.Error(
			"Failed to create user",
			zap.String("command_id", cmd.GetCommandId()),
			zap.String("email", cmd.GetPayload().GetEmail()),
			zap.Error(err),
		)
		return nil
	}

	event := &proto.EventUserCreated{
		Id:       id,
		Name:     cmd.GetPayload().GetName(),
		LastName: cmd.GetPayload().GetLastName(),
		Email:    cmd.GetPayload().GetEmail(),
	}

	channel.Channels.Success(cmd.GetCommandId(), event)

	c.logger.Info(
		"User created",
		zap.String("user_id", id),
	)

	err = c.eventBus.Publish(ctx, event)
	if err != nil {
		c.logger.Error(
			"Failed to publish event.go",
			zap.String("command_id", cmd.GetCommandId()),
			zap.String("email", cmd.GetPayload().GetEmail()),
			zap.Error(err),
		)
		return err
	}

	return nil
}
