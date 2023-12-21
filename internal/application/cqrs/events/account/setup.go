package account

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/startcodextech/goauth/internal/infrastructure/brevo"
	"go.uber.org/zap"
	"os"
)

func SetupHandlers(processor *cqrs.EventProcessor, brevoApi brevo.Brevo, logger *zap.Logger) {
	err := processor.AddHandlers(
		UserCreatedOnCreateUser{
			brevoApi: brevoApi,
			logger:   logger,
		},
		UserCreatedFailedOnCreateUser{},
	)
	if err != nil {
		logger.Error("Failed to add events handler", zap.Error(err))
		os.Exit(1)
	}
}
