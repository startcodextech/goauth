package handlers

import (
	"context"
	"github.com/startcodextech/goauth/users/internal/application"
	"github.com/startcodextech/goauth/users/pb"
	"github.com/startcodextech/goevents/async"
	"github.com/startcodextech/goevents/ddd"
	"github.com/startcodextech/goevents/errorsotel"
	"github.com/startcodextech/goevents/registry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type (
	commandHandlers struct {
		app application.App
	}
)

func NewCommandHandlers(reg registry.Registry, app application.App, replyPublisher async.ReplyPublisher, mws ...async.MessageHandlerMiddleware) async.MessageHandler {
	return async.NewCommandHandler(reg, replyPublisher, commandHandlers{
		app: app,
	}, mws...)
}

func RegisterCommandHandlers(subscriber async.MessageSubscriber, handlers async.MessageHandler) error {
	_, err := subscriber.Subscribe(pb.CommandChannel, handlers, async.MessageFilter{}, async.GroupName("users-commands"))
	return err
}

func (h commandHandlers) HandleCommand(ctx context.Context, cmd ddd.Command) (reply ddd.Reply, err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling command",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled command", trace.WithAttributes(
			attribute.Int64("took-ms", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling command", trace.WithAttributes(
		attribute.String("Command", cmd.CommandName()),
	))

	return nil, nil
}
