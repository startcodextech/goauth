package users

import (
	"context"
	"github.com/startcodextech/goauth/users/internal/constants"
	"github.com/startcodextech/goauth/users/internal/domain"
	"github.com/startcodextech/goauth/users/pb"
	"github.com/startcodextech/goevents/async"
	"github.com/startcodextech/goevents/asyncotel"
	"github.com/startcodextech/goevents/asyncprom"
	"github.com/startcodextech/goevents/ddd"
	"github.com/startcodextech/goevents/depinjection"
	"github.com/startcodextech/goevents/esourcing"
	"github.com/startcodextech/goevents/jetstream"
	"github.com/startcodextech/goevents/registry"
	"github.com/startcodextech/goevents/registry/serdes"
	mongo2 "github.com/startcodextech/goevents/store/mongo"
	mongotel "github.com/startcodextech/goevents/store/mongootel"
	"github.com/startcodextech/goevents/system"
	"github.com/startcodextech/goevents/transmanager"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Module struct {
		Name string
	}
)

func (Module) Startup(ctx context.Context, mono system.Service) (err error) {
	return Root(ctx, mono)
}

func Root(ctx context.Context, svc system.Service) (err error) {
	container := depinjection.New()

	container.AddSingleton(constants.RegistryKey, func(c depinjection.Container) (any, error) {
		reg := registry.New()
		if err := registrations(reg); err != nil {
			return nil, err
		}
		if err := pb.Registrations(reg); err != nil {
			return nil, err
		}
		return reg, nil
	})

	stream := jetstream.NewStream(svc.Config().Nats.Stream, svc.JS(), svc.Logger())
	container.AddSingleton(constants.DomainDispatcherKey, func(c depinjection.Container) (any, error) {
		return ddd.NewEventDispatcher[ddd.Event](), nil
	})

	// database transaction
	container.AddScoped(constants.DatabaseTransactionKey, func(c depinjection.Container) (any, error) {
		session, err := svc.MongoDB().StartSession()
		if err != nil {
			return nil, err
		}

		err = session.StartTransaction()
		if err != nil {
			return nil, err
		}
		return session, nil
	})

	sentCounter := asyncprom.SentMessagesCounter(constants.ServiceName)

	container.AddScoped(constants.MessagePublisherKey, func(c depinjection.Container) (any, error) {
		tx := mongotel.Trace(c.Get(constants.DatabaseTransactionKey).(*mongo.Collection))
		outboxStore := mongo2.NewOutboxStore(tx)

		return async.NewMessagePublisher(
			stream,
			asyncotel.OtelMessageContextInjector(),
			sentCounter,
			transmanager.OutboxPublisher(outboxStore),
		), nil
	})

	container.AddSingleton(constants.MessageSubscriberKey, func(c depinjection.Container) (any, error) {
		return async.NewMessageSubscriber(
			stream,
			asyncotel.OtelMessageContextExtractor(),
			asyncprom.ReceivedMessagesCounter(constants.ServiceName)
		), nil
	})

	container.AddScoped(constants.EventPublisherKey, func(c depinjection.Container) (any, error) {
		return async.NewEventPublisher(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.MessagePublisherKey).(async.MessagePublisher),
		), nil
	})

	container.AddScoped(constants.ReplyPublisherKey, func(c depinjection.Container) (any, error) {
		return async.NewReplyPublisher(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.MessagePublisherKey).(async.MessagePublisher),
		), nil
	})

	container.AddScoped(constants.InboxStoreKey, func(c depinjection.Container) (any, error) {
		tx := mongotel.Trace(c.Get(constants.DatabaseTransactionKey).(*mongo.Collection))
		return mongo2.NewInboxStore(tx), nil
	})

	container.AddScoped(constants.OrdersRepoKey, func(c depinjection.Container) (any, error) {
		tx := mongotel.Trace(c.Get(constants.DatabaseTransactionKey).(*mongo.Collection))
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		return esourcing.NewAggregateRepository[*domain.User](
			domain.UserAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			esourcing.AggregateStoreWithMiddleware(
				mongo2.NewEventStore(constants.EventsTableName, tx, reg),
				mongo2.NewSnapshotStore(tx, reg),
			),
		), nil
	})

	return nil
}

func registrations(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	// Users
	if err = serde.Register(domain.User{}, func(v interface{}) error {
		user := v.(*domain.User)
		user.Aggregate = esourcing.NewAggregate("", domain.UserAggregate)
		return nil
	}); err != nil {
		return err
	}

	// users events
	if err = serde.Register(domain.UserCreated{}); err != nil {
		return err
	}

	// users snapshots
	if err = serde.RegisterKey(domain.UserV1{}.SnapshotName(), domain.UserV1{}); err != nil {
		return err
	}

	return nil
}
