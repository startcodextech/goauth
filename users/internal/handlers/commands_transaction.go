package handlers

import (
	"context"
	"github.com/startcodextech/goauth/users/internal/constants"
	"github.com/startcodextech/goevents/asyncmessages"
	"github.com/startcodextech/goevents/dependencyinjection"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterCommandHandlersTx(container dependencyinjection.Container) error {
	rawMsgHandler := asyncmessages.MessageHandlerFunc(func(ctx context.Context, msg asyncmessages.IncomingMessage) (err error) {
		client := dependencyinjection.Get(ctx, constants.DatabaseTransactionKey).(*mongo.Client)
		session, err := client.StartSession()
		if err != nil {
			return err
		}
		defer session.EndSession(ctx)

		_, err = session.WithTransaction(ctx, func(sessionCtx mongo.SessionContext) (interface{}, error) {
			ctx = mongo.NewSessionContext(ctx, sessionCtx)
			ctx = container.Scoped(ctx)

			err = dependencyinjection.
				Get(ctx, constants.CommandHandlersKey).(asyncmessages.MessageHandler).
				HandleMessage(ctx, msg)
			return nil, err
		})
		return err
	})

	subscriber := container.Get(constants.MessageSubscriberKey).(asyncmessages.MessageSubscriber)

	return RegisterCommandHandlers(subscriber, rawMsgHandler)
}
