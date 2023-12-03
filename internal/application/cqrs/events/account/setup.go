package account

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

func SetupHandlers(processor *cqrs.EventGroupProcessor, logger watermill.LoggerAdapter) {
	err := processor.AddHandlersGroup(
		"events",
		UserCreatedOnCreateUser{},
		UserCreatedFailedOnCreateUser{},
	)
	if err != nil {
		logger.Error("Failed to add events handler", err, nil)
		panic(err)
	}
}
