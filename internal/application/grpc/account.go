package grpc

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/google/uuid"
	"github.com/startcodextech/goauth/internal/application/cqrs/commands"
	"github.com/startcodextech/goauth/proto"
	"github.com/startcodextech/goauth/util/channel"
	"go.uber.org/zap"
	"net/http"
)

type AccountService struct {
	proto.UnimplementedAccountServiceServer
	commandBus *cqrs.CommandBus
	logger     *zap.Logger
}

func NewAccountService(commandBus *cqrs.CommandBus, logger *zap.Logger) *AccountService {
	return &AccountService{
		commandBus: commandBus,
		logger:     logger,
	}
}

func (s *AccountService) CreateUser(ctx context.Context, request *proto.CreateUserRequest) (*proto.ResponseWithString, error) {
	result := &proto.ResponseWithString{
		Status: http.StatusCreated,
	}

	cmd := &proto.CommandCreateUser{
		CommandId: uuid.New().String(),
		Payload:   request.GetUser(),
	}

	response := make(chan channel.ResultChannel, 1)
	channel.AddChannel(cmd.GetCommandId(), response)

	err := commands.Publish(ctx, s.commandBus, cmd.GetCommandId(), cmd, s.logger)
	if err != nil {
		result.Status = http.StatusInternalServerError
		result.Error = err.Error()
		return result, nil
	}

	err = channel.GetResult(response, channel.ResultCallback{
		CorrelationID: cmd.GetCommandId(),
		OnSuccess: func(i interface{}) {
			data := i.(*proto.EventUserCreated)
			result.Data = data.GetId()
		},
		OnFailed: func(err *proto.EventError) {
			result.Error = err.GetError()
			result.Status = http.StatusBadRequest
		},
		Logger: s.logger,
	})
	if err != nil {
		result.Status = http.StatusInternalServerError
		result.Error = err.Error()
		return result, nil
	}

	return result, err
}
