package handlers

import (
	"context"
	"github.com/startcodextech/goauth/users/internal/constants"
	"github.com/startcodextech/goevents/async"
	"github.com/startcodextech/goevents/depinjection"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterCommandHandlersTx(container depinjection.Container) error {
	rawMsgHandler := async.MessageHandlerFunc(func(ctx context.Context, msg async.IncomingMessage) (err error) {
		client := depinjection.Get(ctx, constants.DatabaseTransactionKey).(*mongo.Client)
		session, err := client.StartSession()
		if err != nil {
			return err
		}
		defer session.EndSession(ctx)

		_, err = session.WithTransaction(ctx, func(sessionCtx mongo.SessionContext) (interface{}, error) {
			ctx = mongo.NewSessionContext(ctx, sessionCtx)
			ctx = container.Scoped(ctx)

			err = depinjection.
				Get(ctx, constants.CommandHandlersKey).(async.MessageHandler).
				HandleMessage(ctx, msg)
			return nil, err
		})
		return err
	})

	subscriber := container.Get(constants.MessageSubscriberKey).(async.MessageSubscriber)

	return RegisterCommandHandlers(subscriber, rawMsgHandler)
}
