package cqrs

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"time"
)

// NewCqrsMarshaler returns a new cqrs marshaler.
// It is used to marshal/unmarshal commands and events.
func NewCqrsMarshaler() cqrs.CommandEventMarshaler {
	return ProtoMarshal{}
}

// NewCommandBus creates a new command bus.
// It is used to send commands to the command handler.
func NewCommandBus(commandPublisher message.Publisher, marshaler cqrs.CommandEventMarshaler, logger watermill.LoggerAdapter) *cqrs.CommandBus {
	bus, err := cqrs.NewCommandBusWithConfig(commandPublisher, cqrs.CommandBusConfig{
		GeneratePublishTopic: func(params cqrs.CommandBusGeneratePublishTopicParams) (string, error) {
			return params.CommandName, nil
		},
		OnSend: func(params cqrs.CommandBusOnSendParams) error {
			logger.Info("Sending command", watermill.LogFields{
				"command_name": params.CommandName,
			})

			params.Message.Metadata.Set("sent_at", time.Now().String())

			return nil
		},
		Marshaler: marshaler,
		Logger:    logger,
	})
	if err != nil {
		logger.Error("Failed to create command bus", err, nil)
		panic(err)
	}

	return bus
}

// NewCqrsRouter creates a new cqrs router.
// It is used to route commands to the command handler.
func NewCqrsRouter(logger watermill.LoggerAdapter) *message.Router {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		logger.Error("Failed to create router", err, nil)
		panic(err)
	}

	retryMiddleware := middleware.Retry{
		MaxRetries:      0,
		InitialInterval: 2 * time.Second,
		Logger:          logger,
	}

	router.AddMiddleware(retryMiddleware.Middleware)
	router.AddMiddleware(middleware.Recoverer)

	return router
}

// NewCommandProcessor creates a new command processor.
// It is used to handle commands.
func NewCommandProcessor(commandSubscriber message.Subscriber, router *message.Router, marshaler cqrs.CommandEventMarshaler, logger watermill.LoggerAdapter) *cqrs.CommandProcessor {
	cmdProcessor, err := cqrs.NewCommandProcessorWithConfig(
		router,
		cqrs.CommandProcessorConfig{
			GenerateSubscribeTopic: func(params cqrs.CommandProcessorGenerateSubscribeTopicParams) (string, error) {
				return params.CommandName, nil
			},
			SubscriberConstructor: func(params cqrs.CommandProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return commandSubscriber, nil
			},
			OnHandle: func(params cqrs.CommandProcessorOnHandleParams) error {
				start := time.Now()

				err := params.Handler.Handle(params.Message.Context(), params.Command)

				logger.Info("Command handled", watermill.LogFields{
					"command_name": params.CommandName,
					"duration":     time.Since(start),
					"err":          err,
				})

				params.Message.Ack()

				return err
			},
			Marshaler: marshaler,
			Logger:    logger,
		},
	)
	if err != nil {
		logger.Error("Failed to create command processor", err, nil)
		panic(err)
	}

	return cmdProcessor
}

// NewEventBus creates a new event.go bus.
// It is used to publish events.
func NewEventBus(eventPublisher message.Publisher, marshaler cqrs.CommandEventMarshaler, logger watermill.LoggerAdapter) *cqrs.EventBus {
	bus, err := cqrs.NewEventBusWithConfig(eventPublisher, cqrs.EventBusConfig{
		GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error) {
			return params.EventName, nil
		},
		OnPublish: func(params cqrs.OnEventSendParams) error {
			logger.Info("Publishing event.go", watermill.LogFields{
				"event_name": params.EventName,
			})

			params.Message.Metadata.Set("published_at", time.Now().String())

			return nil
		},
		Marshaler: marshaler,
		Logger:    logger,
	})
	if err != nil {
		logger.Error("Failed to create event.go bus", err, nil)
		panic(err)
	}

	return bus
}

// NewEventProcessor creates a new event.go processor.
// It is used to handle events.
func NewEventProcessor(eventSubscriber message.Subscriber, router *message.Router, marshaler cqrs.CommandEventMarshaler, logger watermill.LoggerAdapter) *cqrs.EventProcessor {
	processor, err := cqrs.NewEventProcessorWithConfig(
		router,
		cqrs.EventProcessorConfig{
			GenerateSubscribeTopic: func(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
				return params.EventName, nil
			},
			SubscriberConstructor: func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return eventSubscriber, nil
			},
			OnHandle: func(params cqrs.EventProcessorOnHandleParams) error {
				start := time.Now()

				err := params.Handler.Handle(params.Message.Context(), params.Event)

				logger.Info("Event handled", watermill.LogFields{
					"event_name": params.EventName,
					"duration":   time.Since(start),
					"err":        err,
				})

				return err
			},
			Marshaler: marshaler,
			Logger:    logger,
		},
	)
	if err != nil {
		logger.Error("Failed to create event.go processor", err, nil)
		panic(err)
	}

	return processor
}
