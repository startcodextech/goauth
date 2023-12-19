package main

import (
	"context"
	_ "github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/startcodextech/goauth/internal/application/cqrs"
	"github.com/startcodextech/goauth/internal/application/cqrs/commands"
	"github.com/startcodextech/goauth/internal/application/cqrs/events"
	"github.com/startcodextech/goauth/internal/application/grpc"
	"github.com/startcodextech/goauth/internal/application/http"
	"github.com/startcodextech/goauth/internal/application/services"
	"github.com/startcodextech/goauth/internal/infrastructure/messaging/gochannel"
	"github.com/startcodextech/goauth/internal/infrastructure/persistence/mongodb"
	"github.com/startcodextech/goauth/proto"
	"github.com/startcodextech/goauth/util/log"
	"go.uber.org/zap"
	"os"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	zapLogger, _ := zap.NewProduction()
	defer func(zapLogger *zap.Logger) {
		err := zapLogger.Sync()
		if err != nil {
			zapLogger.Error("An error occurred while synchronizing the log")
		}
	}(zapLogger)
	logger := log.NewLogger(zapLogger)

	mongo := mongodb.New(ctx, os.Getenv("DB_NAME"), logger)
	defer mongo.Disconnect(ctx)()
	mongo.Ping(ctx)

	pubSub := gochannel.New(logger)
	defer func() {
		if err := pubSub.Close(); err != nil {
			logger.Error("Failed to close pubSub", err, nil)
		}
	}()

	svcs := services.New(ctx, mongo)

	cqrsMarshaler := cqrs.NewCqrsMarshaler()
	cqrsRouter := cqrs.NewCqrsRouter(logger)
	commandBus := cqrs.NewCommandBus(pubSub, cqrsMarshaler, logger)
	commandProcessor := cqrs.NewCommandProcessor(pubSub, cqrsRouter, cqrsMarshaler, logger)
	eventBus := cqrs.NewEventBus(pubSub, cqrsMarshaler, logger)
	eventProcessor := cqrs.NewEventProcessor(pubSub, cqrsRouter, cqrsMarshaler, logger)

	commands.RunHandlers(commandProcessor, eventBus, svcs, zapLogger)
	events.RunHandlers(eventProcessor, eventBus, svcs, zapLogger)

	eventsChannel := make(chan events.EventData)

	messages, err := pubSub.Subscribe(ctx, "proto.EventUserCreatedFailed")
	if err == nil {
		go func(messages <-chan *message.Message) {
			for msg := range messages {

				var data proto.EventUserCreatedFailed

				err := cqrsMarshaler.Unmarshal(msg, &data)
				if err == nil {
					msg.Ack()
					eventData := map[string]interface{}{
						"email": data.Email,
						"error": data.Error,
					}
					eventsChannel <- eventData
				} else {
				}
			}
		}(messages)
	}

	httpServer := http.New(zapLogger)

	rpcServer, err := grpc.New(ctx, httpServer.App(), commandBus, eventsChannel, logger)
	if err != nil {
		zapLogger.Error("", zap.Error(err))
		os.Exit(1)
	}

	rpcServer.Start()
	httpServer.Start()

	if err := cqrsRouter.Run(ctx); err != nil {
		logger.Error("Failed to run cqrs router", err, nil)
		panic(err)
	}

}
