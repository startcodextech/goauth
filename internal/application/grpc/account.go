package grpc

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/google/uuid"
	"github.com/startcodextech/goauth/proto"
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

	err := s.commandBus.Send(ctx, cmd)
	if err != nil {
		s.logger.Error(
			"An error occurred while sending the command",
			zap.String("email", request.GetUser().GetEmail()),
			zap.Error(err),
		)
		result.Status = http.StatusInternalServerError
		result.Error = err.Error()
		return result, nil
	}

	result.Data = cmd.GetCommandId()

	return result, err
}
