package grpc

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/startcodextech/goauth/internal/application/cqrs/events"
	"github.com/startcodextech/goauth/proto"
	"net/http"
	"time"
)

type AccountService struct {
	proto.UnimplementedAccountServiceServer
	commandBus   *cqrs.CommandBus
	logger       watermill.LoggerAdapter
	eventChannel chan events.EventData
}

func NewAccountService(commandBus *cqrs.CommandBus, logger watermill.LoggerAdapter, dataChanel chan events.EventData) *AccountService {
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
		s.logger.Error("", err, watermill.LogFields{
			"email": request.GetUser().GetEmail(),
		})
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
