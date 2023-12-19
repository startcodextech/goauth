package account

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

func SetupHandlers(processor *cqrs.EventProcessor, logger watermill.LoggerAdapter) {
	err := processor.AddHandlers(
		UserCreatedOnCreateUser{},
		UserCreatedFailedOnCreateUser{},
	)
	if err != nil {
		logger.Error("Failed to add events handler", err, nil)
		panic(err)
	}
}
