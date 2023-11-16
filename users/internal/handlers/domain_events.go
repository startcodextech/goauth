package handlers

import (
	"context"
	"github.com/startcodextech/goauth/users/internal/domain"
	"github.com/startcodextech/goauth/users/pb"
	"github.com/startcodextech/goevents/async"
	"github.com/startcodextech/goevents/ddd"
	"github.com/startcodextech/goevents/errorsotel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type (
	domainHandlers[T ddd.Event] struct {
		publisher async.EventPublisher
	}
)

func NewDomainEventHandlers(publisher async.EventPublisher) ddd.EventHandler[ddd.Event] {
	return domainHandlers[ddd.Event]{publisher: publisher}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlerrs ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlerrs,
		domain.UserCreatedEvent,
	)
}

func (h domainHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling domain event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled domain event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handled domain event", trace.WithAttributes(
		attribute.String("event", event.EventName()),
	))

	switch event.EventName() {
	case domain.UserCreatedEvent:
		return h.onUserCreated(ctx, event)
	}

	return nil
}

func (h domainHandlers[T]) onUserCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.User)
	return h.publisher.Publish(ctx, pb.UserAggregateChannel,
		ddd.NewEvent(pb.UserCreatedEvent, &pb.UserCreated{
			Id:       payload.ID(),
			Email:    payload.Email,
			Phone:    payload.Phone,
			Password: payload.PasswordHash,
			Name:     payload.Name,
			LastName: payload.LastName,
			Enable:   payload.Enabled,
		}),
	)
}
