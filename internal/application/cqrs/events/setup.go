package events

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/startcodextech/goauth/internal/application/cqrs/events/account"
	"github.com/startcodextech/goauth/internal/application/cqrs/events/types"
	"github.com/startcodextech/goauth/internal/application/services"
	"go.uber.org/zap"
)

func RunHandlers(processor *cqrs.EventProcessor, eventBus *cqrs.EventBus, services services.Services, logger *zap.Logger) {
	account.SetupHandlers(processor, logger)
}

func RunSubscriber(ctx context.Context, subscriber message.Subscriber, marshaler cqrs.CommandEventMarshaler, eventsChannel chan types.EventData, logger *zap.Logger) {
	account.SetupSubscriber(ctx, subscriber, marshaler, eventsChannel, logger)
}
