package account

import (
	"context"
	"github.com/startcodextech/goauth/internal/infrastructure/brevo"
	"github.com/startcodextech/goauth/proto"
	"go.uber.org/zap"
	"log"
)

type (
	UserCreatedOnCreateUser struct {
		brevoApi brevo.Brevo
		logger   *zap.Logger
	}

	UserCreatedFailedOnCreateUser struct {
	}
)

func (e UserCreatedOnCreateUser) HandlerName() string {
	return "UserCreatedOnCreateUser"
}

func (UserCreatedOnCreateUser) NewEvent() interface{} {
	return &proto.EventUserCreated{}
}

func (e UserCreatedOnCreateUser) Handle(ctx context.Context, event interface{}) error {
	eventMsg, ok := event.(*proto.EventUserCreated)
	if !ok {
		return nil
	}

	params := map[string]interface{}{
		"name": eventMsg.Name,
		"url":  "https://startcodex.com/user/verify?token=sdkljfh8eir32werwnkwjchrewihcrwejndkweuhywieucfw",
	}

	sendTo := []brevo.SendTo{
		{
			Email: eventMsg.Email,
			Name:  eventMsg.Name,
		},
	}

	err := e.brevoApi.SendTemplateEmail(ctx, true, 1, sendTo, params)
	if err != nil {
		return err
	}

	return nil
}

func (e UserCreatedFailedOnCreateUser) HandlerName() string {
	return "UserCreatedFailedOnCreateUser"
}

func (UserCreatedFailedOnCreateUser) NewEvent() interface{} {
	return &proto.EventError{}
}

func (e UserCreatedFailedOnCreateUser) Handle(ctx context.Context, event interface{}) error {
	eventMsg, ok := event.(*proto.EventError)
	if !ok {
		log.Printf("Se recibió un tipo de evento inesperado: %T", event)
		return nil
	}

	log.Println("------------------------------------")
	log.Printf("UserCreatedFailedOnCreateUser: %s", eventMsg)
	log.Println("------------------------------------")

	return nil
}
