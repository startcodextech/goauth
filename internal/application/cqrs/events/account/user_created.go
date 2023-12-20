package account

import (
	"context"
	"github.com/startcodextech/goauth/proto"
	"log"
)

type (
	UserCreatedOnCreateUser struct {
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
		log.Printf("Se recibió un tipo de evento inesperado: %T", event)
		return nil
	}

	log.Println("------------------------------------")
	log.Printf("UserCreatedOnCreateUser: %s", eventMsg)
	log.Println("------------------------------------")

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
