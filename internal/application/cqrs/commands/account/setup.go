package account

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/startcodextech/goauth/internal/application/services/account"
)

func SetupHandlers(processor *cqrs.CommandProcessor, eventBus *cqrs.EventBus, services account.AccountServices, logger watermill.LoggerAdapter) {
	err := processor.AddHandlers(
		CreateUserHandler{
			eventBus: eventBus,
			service:  services.User(),
			logger:   logger,
		},
	)

	if err != nil {
		logger.Error("Failed to add account command handler", err, nil)
		panic(err)
	}
}
