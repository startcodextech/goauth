package commands

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/startcodextech/goauth/internal/application/cqrs/commands/account"
	"github.com/startcodextech/goauth/internal/application/services"
	"go.uber.org/zap"
)

func RunHandlers(processor *cqrs.CommandProcessor, eventBus *cqrs.EventBus, services services.Services, logger *zap.Logger) {
	account.SetupHandlers(processor, eventBus, services.Account(), logger)
}
