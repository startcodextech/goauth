package grpc

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/startcodextech/goauth/internal/application/cqrs/events/types"
	"github.com/startcodextech/goauth/proto"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type AccountService struct {
	proto.UnimplementedAccountServiceServer
	commandBus   *cqrs.CommandBus
	logger       *zap.Logger
	eventChannel chan types.EventData
}

func NewAccountService(commandBus *cqrs.CommandBus, logger *zap.Logger, dataChanel chan types.EventData) *AccountService {
	return &AccountService{
		commandBus:   commandBus,
		logger:       logger,
		eventChannel: dataChanel,
	}
}

func (s *AccountService) CreateUser(ctx context.Context, request *proto.CreateUserRequest) (*proto.StandardResponseWithString, error) {

	reading := true

	result := &proto.StandardResponseWithString{
		Status: http.StatusCreated,
	}

	err := s.commandBus.Send(ctx, request.GetUser())
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

	for reading {
		select {
		case eventData := <-s.eventChannel:
			if eventData["email"] == request.GetUser().GetEmail() {
				reading = false
				if err, ok := eventData["error"]; ok {
					result.Error = err.(string)
					result.Status = http.StatusBadRequest
				}
				id, ok := eventData["id"]
				if !ok {
					result.Status = http.StatusBadRequest
				} else {
					result.Data = id.(string)
				}
			}
		case <-time.After(30 * time.Second):
			reading = false
		}
	}

	return result, nil
}
