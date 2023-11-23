package http

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/modernice/goes/command"
	"github.com/startcodextech/goauth/internal/domain/account"
)

type AccountRest struct {
	ctx    context.Context
	router fiber.Router
	bus    command.Bus
}

func NewAccountRest(ctx context.Context, router fiber.Router, bus command.Bus) AccountRest {
	return AccountRest{
		ctx:    ctx,
		router: router,
		bus:    bus,
	}
}

func (r AccountRest) Setup() {

	v1 := r.router.Group("/v1/account")

	v1.Post("/user", r.create)
}

func (r AccountRest) create(ctx *fiber.Ctx) error {
	var body account.UserCreateDto

	if err := ctx.BodyParser(&body); err != nil {
		return err
	}

	id := uuid.New()

	cmd := account.CreateUser(id, body)
	if err := r.bus.Dispatch(r.ctx, cmd.Any()); err != nil {
		return err
	}

	return nil
}
