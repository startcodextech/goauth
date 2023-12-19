package account

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/startcodextech/goauth/internal/application/services/account"
	"go.uber.org/zap"
	"os"
)

func SetupHandlers(processor *cqrs.CommandProcessor, eventBus *cqrs.EventBus, services account.AccountServices, logger *zap.Logger) {
	err := processor.AddHandlers(
		CreateUserHandler{
			eventBus: eventBus,
			service:  services.User(),
			logger:   logger,
		},
	)

	if err != nil {
		logger.Error("Failed to add account command handler", zap.Error(err))
		os.Exit(1)
	}
}
