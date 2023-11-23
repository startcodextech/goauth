package application

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/modernice/goes/aggregate/repository"
	"github.com/modernice/goes/backend/mongo"
	"github.com/modernice/goes/backend/nats"
	"github.com/modernice/goes/codec"
	"github.com/modernice/goes/command"
	"github.com/modernice/goes/command/cmdbus"
	"github.com/modernice/goes/event"
	"github.com/modernice/goes/event/eventstore"
	"github.com/startcodextech/goauth/internal/domain/account"
	"github.com/startcodextech/goauth/internal/interfaces/http"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Setup struct {
	http            *fiber.App
	ctx             context.Context
	cancelCtx       context.CancelFunc
	eventBus        event.Bus
	eventDisconnect func()
	commandBus      command.Bus
	store           event.Store
	eventReg        *codec.Registry
	cmdReg          *codec.Registry
	serviceName     string
}

func New(serviceName string) Setup {

	httpApp := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	httpApp.Use(cache.New())
	httpApp.Use(compress.New())
	httpApp.Use(cors.New())
	//httpApp.Use(csrf.New())
	httpApp.Use(logger.New())
	httpApp.Use(recover.New())

	ctx, cancelCtx := newContext()

	eventReg := event.NewRegistry()
	bus, disconnect := newEventBus(ctx, eventReg)

	return Setup{
		http:            httpApp,
		ctx:             ctx,
		cancelCtx:       cancelCtx,
		eventReg:        eventReg,
		cmdReg:          command.NewRegistry(),
		eventBus:        bus,
		eventDisconnect: disconnect,
	}
}

func newContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
}

func newEventBus(ctx context.Context, enc codec.Encoding) (event.Bus, func()) {
	bus := nats.NewEventBus(enc)

	return bus, func() {
		log.Printf("Disconnecting from NATS ...")

		if err := bus.Disconnect(ctx); err != nil {
			log.Panicf("Failed to disconnect from NATS: %v", err)
		}
	}
}

func (s *Setup) Cancel() {
	s.cancelCtx()
}

func (s *Setup) Disconnect() {
	s.eventDisconnect()
}

func (s *Setup) Context() context.Context {
	return s.ctx
}

func (s *Setup) CommandBus() command.Bus {
	return s.commandBus
}

func (s *Setup) Events() {

	log.Printf("Setting up events ...")

	account.RegisterEvents(s.eventReg)

	s.store = eventstore.WithBus(mongo.NewEventStore(s.eventReg), s.eventBus)
}

func (s *Setup) Commands() {
	log.Printf("Setting up commands ...")

	account.RegisterCommands(s.cmdReg)

	cmdbus.RegisterEvents(s.eventReg)

	s.commandBus = cmdbus.New[int](s.cmdReg, s.eventBus)
}

func (s *Setup) Aggregates() *repository.Repository {
	log.Printf("Setting up aggregates")

	return repository.New(s.store)
}

func (s *Setup) Rest() {

	api := s.http.Group("/api")

	http.NewAccountRest(s.ctx, api, s.commandBus).Setup()

	s.http.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.Send([]byte("Service up"))
	})

	go s.http.Listen(":8080")
}
