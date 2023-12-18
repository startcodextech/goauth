package grpc

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/startcodextech/goauth/proto"
	"log"
	"net/http"
)

type AccountService struct {
	proto.UnimplementedAccountServiceServer
	commandBus      *cqrs.CommandBus
	eventSubscriber message.Subscriber
	logger          watermill.LoggerAdapter
}

func NewAccountService(commandBus *cqrs.CommandBus, eSubscriber message.Subscriber, logger watermill.LoggerAdapter) *AccountService {
	return &AccountService{
		commandBus:      commandBus,
		eventSubscriber: eSubscriber,
		logger:          logger,
	}
}

func (s *AccountService) CreateUser(ctx context.Context, request *proto.CreateUserRequest) (*proto.StandardResponseWithString, error) {
	s.logger.Info("crating user", watermill.LogFields{
		"email": request.GetUser().GetEmail(),
	})

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

	messages, err := s.eventSubscriber.Subscribe(ctx, "events")
	if err != nil {
		s.logger.Error("", err, watermill.LogFields{
			"email": request.GetUser().GetEmail(),
		})
		result.Status = http.StatusInternalServerError
		result.Error = err.Error()
		return result, nil
	}

	go func(messages <-chan *message.Message) {
		for msg := range messages {
			log.Printf("received message: %s", msg.Payload)
		}
	}(messages)

	return result, nil
}
