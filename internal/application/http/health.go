package http

import "github.com/gofiber/fiber/v2"

func healthEndPoint() (string, func(ctx *fiber.Ctx) error) {
	return "health", func(ctx *fiber.Ctx) error {
		return ctx.SendString("OK")
	}
}
