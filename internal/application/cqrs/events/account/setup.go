package account

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/startcodextech/goauth/internal/application/cqrs/events/types"
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

func SetupSubscriber(ctx context.Context, subscriber message.Subscriber, marshaler cqrs.CommandEventMarshaler, eventsChannel chan types.EventData, logger *zap.Logger) {
	subscriberUserCreatedSuccess(ctx, subscriber, marshaler, eventsChannel, logger)
	subscriberUserCreatedFailed(ctx, subscriber, marshaler, eventsChannel, logger)
}
