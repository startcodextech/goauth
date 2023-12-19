package events

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/startcodextech/goauth/internal/application/cqrs/events/account"
	"github.com/startcodextech/goauth/internal/application/services"
	"go.uber.org/zap"
)

func RunHandlers(processor *cqrs.EventProcessor, eventBus *cqrs.EventBus, services services.Services, logger *zap.Logger) {
	account.SetupHandlers(processor, logger)
}
