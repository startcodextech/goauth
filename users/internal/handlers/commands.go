package handlers

import (
	"github.com/startcodextech/goauth/users/internal/application"
	"github.com/startcodextech/goevents/asyncmessages"
	"github.com/startcodextech/goevents/registry"
)

type (
	commandHandlers struct {
		app application.App
	}
)

func NewCommandHandlers(reg registry.Registry, app application.App, replyPublisher asyncmessages.ReplyPublisher, mws ...asyncmessages.MessageHandlerMiddleware) asyncmessages.MessageHandler {
	return asyncmessages.NewCommandHandler(reg, replyPublisher, commandHandlers{
		app: app,
	}, mws)
}

func RegisterCommandHandlers(subscriber asyncmessages.MessageSubscriber, handlers asyncmessages.MessageHandler) error {
	_, err := subscriber.Subscribe(CommandChannel)
}
