package events

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/startcodextech/goauth/internal/application/cqrs/events/account"
	"github.com/startcodextech/goauth/internal/application/services"
	"github.com/startcodextech/goauth/internal/infrastructure/brevo"
	"go.uber.org/zap"
)

func RunHandlers(
	processor *cqrs.EventProcessor,
	eventBus *cqrs.EventBus,
	services services.Services,
	brevoApi *brevo.Brevo,
	logger *zap.Logger,
) {
	account.SetupHandlers(processor, brevoApi, logger)
}
