package events

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/startcodextech/goauth/internal/application/cqrs/events/account"
	"github.com/startcodextech/goauth/internal/application/services"
)

func RunHandlers(processor *cqrs.EventGroupProcessor, eventBus *cqrs.EventBus, services services.Services, logger watermill.LoggerAdapter) {
	account.SetupHandlers(processor, logger)
}
