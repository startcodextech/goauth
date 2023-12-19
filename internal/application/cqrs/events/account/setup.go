package account

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"go.uber.org/zap"
	"os"
)

func SetupHandlers(processor *cqrs.EventProcessor, logger *zap.Logger) {
	err := processor.AddHandlers(
		UserCreatedOnCreateUser{},
		UserCreatedFailedOnCreateUser{},
	)
	if err != nil {
		logger.Error("Failed to add events handler", zap.Error(err))
		os.Exit(1)
	}
}
