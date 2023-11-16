package handlers

import (
	"context"
	"github.com/startcodextech/goauth/users/internal/constants"
	"github.com/startcodextech/goevents/ddd"
	"github.com/startcodextech/goevents/depinjection"
)

func RegisterDomainEventHandlersTx(container depinjection.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		domainHandlers := depinjection.Get(ctx, constants.DomainEventHandlersKey).(ddd.EventHandler[ddd.Event])

		return domainHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])

	RegisterDomainEventHandlers(subscriber, handlers)
}
