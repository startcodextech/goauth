package http

import (
	"errors"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

const (
	httpAddress = ":8000"
)

func New() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: defaultErrorHandler(),
	})

	app.Use(cors.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	app.Use(recover.New())

	return app
}

func defaultErrorHandler() func(ctx *fiber.Ctx, err error) error {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError

		// Retrieve the custom status code if it's a *fiber.Error
		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
		}

		// Send custom error
		err = ctx.Status(code).JSON(map[string]interface{}{
			"status": code,
			"error":  err.Error(),
		})
		if err != nil {
			// In case the SendFile fails
			return ctx.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{
				"code":  fiber.StatusInternalServerError,
				"error": "Internal server error",
			})
		}

		// Return from handler
		return nil
	}
}

func Start(app *fiber.App, logger watermill.LoggerAdapter) {

	setup(app)

	err := app.Listen(httpAddress)
	if err != nil {
		logger.Error("", err, nil)
		panic(err)
	}

}

func setup(app *fiber.App) {
	app.Get("health", func(ctx *fiber.Ctx) error {
		return ctx.SendString("OK")
	})
}
