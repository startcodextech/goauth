package constants

const (
	// ServiceName the name of this module/service
	ServiceName = "users"

	// Dependency Injection Key
	RegistryKey            = "registry"
	DomainDispatcherKey    = "domainDispatcher"
	DatabaseTransactionKey = "tx"
	MessagePublisherKey    = "messagePublisher"
	MessageSubscriberKey   = "messageSubscriber"
	EventPublisherKey      = "eventPublisher"
	ReplyPublisherKey      = "replyPublisher"
	InboxStoreKey          = "inboxStore"
	DomainEventHandlersKey = "domainEventHandlers"
	CommandHandlersKey     = "commandHandlers"

	OrdersRepoKey = "usersRepo"

	// Repository tables names
	OutboxTableName    = ServiceName + ".outbox"
	InboxTableName     = ServiceName + ".inbox"
	EventsTableName    = ServiceName + ".events"
	SnapshotsTableName = ServiceName + ".snapshots"
)
