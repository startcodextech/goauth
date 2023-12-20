package http

import (
	"errors"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

const (
	httpAddress = ":80"
)

type Server struct {
	logger  *zap.Logger
	server  *fiber.App
	address string
}

func New(logger *zap.Logger) *Server {

	address := os.Getenv("PORT")
	if len(address) == 0 {
		address = httpAddress
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: defaultErrorHandler(),
	})

	app.Use(cors.New())
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
	}))
	app.Use(recover.New())

	basePath, err := os.Getwd()
	if err == nil {
		app.Use(swagger.New(
			swagger.Config{
				BasePath: basePath,
				FilePath: "./swagger.json",
				Path:     "swagger",
				Title:    "API Documentation",
			}),
		)
	}

	return &Server{
		logger:  logger,
		server:  app,
		address: address,
	}
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

func (s *Server) Start() {
	s.setup()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stopChan

		s.logger.Info("Initiating controlled server HTTP shutdown...")
		if err := s.server.Shutdown(); err != nil {
			s.logger.Error("Error shutting down HTTP server", zap.Error(err))
			return
		}
		s.logger.Info("Server HTTP shut down successfully")
	}()

	go func() {
		err := s.server.Listen(httpAddress)
		if err != nil {
			s.logger.Error("An error occurred starting the http server", zap.Error(err))
		}
	}()
}

func (s *Server) setup() {
	s.server.Get(healthEndPoint())
}

func (s *Server) App() *fiber.App {
	return s.server
}
