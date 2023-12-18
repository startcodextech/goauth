package main

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/startcodextech/goauth/internal/application/cqrs"
	"github.com/startcodextech/goauth/internal/application/cqrs/commands"
	"github.com/startcodextech/goauth/internal/application/cqrs/events"
	"github.com/startcodextech/goauth/internal/application/services"
	"github.com/startcodextech/goauth/internal/infrastructure/grpc"
	"github.com/startcodextech/goauth/internal/infrastructure/messaging/gochannel"
	"github.com/startcodextech/goauth/internal/infrastructure/persistence/mongodb"
	"os"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	logger := watermill.NewStdLogger(false, true)

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

	commands.RunHandlers(commandProcessor, eventBus, svcs, logger)
	events.RunHandlers(eventProcessor, eventBus, svcs, logger)

	grpc.Start(ctx, commandBus, pubSub, logger)

	if err := cqrsRouter.Run(ctx); err != nil {
		logger.Error("Failed to run cqrs router", err, nil)
		panic(err)
	}

}
