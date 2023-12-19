package account

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/startcodextech/goauth/internal/application/cqrs/events/types"
	"github.com/startcodextech/goauth/proto"
	"go.uber.org/zap"
)

func subscriberUserCreatedSuccess(ctx context.Context, subscriber message.Subscriber, marshaler cqrs.CommandEventMarshaler, eventsChannel chan types.EventData, logger *zap.Logger) {
	messages, err := subscriber.Subscribe(ctx, "proto.EventUserCreated")
	if err == nil {
		go func(messages <-chan *message.Message) {
			for msg := range messages {
				var data proto.EventUserCreated

				err := marshaler.Unmarshal(msg, &data)
				if err != nil {
					logger.Error("An error occurred while Unmarshal EventUserCreated", zap.Error(err))
					msg.Nack()
				} else {
					msg.Ack()
					eventData := map[string]interface{}{
						"email": data.Email,
						"id":    data.Id,
					}
					eventsChannel <- eventData
				}
			}
		}(messages)
	}
}

func subscriberUserCreatedFailed(ctx context.Context, subscriber message.Subscriber, marshaler cqrs.CommandEventMarshaler, eventsChannel chan types.EventData, logger *zap.Logger) {
	messages, err := subscriber.Subscribe(ctx, "proto.EventUserCreatedFailed")
	if err == nil {
		go func(messages <-chan *message.Message) {
			for msg := range messages {
				var data proto.EventUserCreatedFailed

				err := marshaler.Unmarshal(msg, &data)
				if err != nil {
					logger.Error("An error occurred while Unmarshal EventUserCreatedFailed", zap.Error(err))
					msg.Nack()
				} else {
					msg.Ack()
					eventData := map[string]interface{}{
						"email": data.Email,
						"error": data.Error,
					}
					eventsChannel <- eventData
				}
			}
		}(messages)
	}
}
