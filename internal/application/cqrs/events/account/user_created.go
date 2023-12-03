package account

import (
	"context"
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
	return &UserCreated{}
}

func (e UserCreatedOnCreateUser) Handle(ctx context.Context, event interface{}) error {
	eventMsg := event.(*UserCreated)

	log.Printf("UserCreatedOnCreateUser: %s", eventMsg)

	return nil
}

func (e UserCreatedFailedOnCreateUser) HandlerName() string {
	return "UserCreatedFailedOnCreateUser"
}

func (UserCreatedFailedOnCreateUser) NewEvent() interface{} {
	return &UserCreatedFailed{}
}

func (e UserCreatedFailedOnCreateUser) Handle(ctx context.Context, event interface{}) error {
	eventMsg := event.(*UserCreatedFailed)

	log.Printf("UserCreatedFailedOnCreateUser: %s", eventMsg)

	return nil
}
