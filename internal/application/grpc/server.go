package grpc

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/startcodextech/goauth/internal/application/cqrs/events/types"
	"github.com/startcodextech/goauth/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	grpcNetwork = "tcp"
	grpcAddress = "0.0.0.0:9090"
)

var (
	allowedHeaders = map[string]struct{}{
		"x-request-id": {},
	}
)

type Server struct {
	ctx         context.Context
	commandBus  *cqrs.CommandBus
	dataChannel chan types.EventData
	logger      *zap.Logger
	http        *fiber.App
	listen      net.Listener
	server      *grpc.Server
}

func New(ctx context.Context, serverHTTP *fiber.App, commandBus *cqrs.CommandBus, dataChanel chan types.EventData, logger *zap.Logger) (*Server, error) {

	address := os.Getenv("PORT_RPC")
	if len(address) == 0 {
		address = grpcAddress
	}

	listen, err := net.Listen(grpcNetwork, address)
	if err != nil {
		return nil, err
	}

	return &Server{
		ctx:         ctx,
		commandBus:  commandBus,
		dataChannel: dataChanel,
		logger:      logger,
		http:        serverHTTP,
		listen:      listen,
		server:      grpc.NewServer(),
	}, nil
}

func (s *Server) Start() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stopChan

		s.logger.Info("Initiating controlled server gRPC shutdown...")
		s.server.GracefulStop()
		s.logger.Info("Server gRPC shut down successfully")
	}()

	go func() {
		err := s.startGRPCServer()
		if err != nil {
			s.logger.Error("An error occurred starting the grpc server", zap.Error(err))
			os.Exit(1)
		}
	}()

	s.configureMux()
}

func (s *Server) startGRPCServer() error {
	accountService := NewAccountService(s.commandBus, s.logger, s.dataChannel)
	proto.RegisterAccountServiceServer(s.server, accountService)

	err := s.server.Serve(s.listen)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) configureMux() {
	mux := s.configureHTTPGateway()

	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := proto.RegisterAccountServiceHandlerFromEndpoint(s.ctx, mux, s.listen.Addr().String(), dialOptions)
	if err != nil {
		s.logger.Error(
			"An error occurred while registering the endpoint in the account service handle",
			zap.Error(err),
		)
		os.Exit(1)
	}

	s.http.Group("api/v1/*", adaptor.HTTPHandler(mux))
}

func (s *Server) configureHTTPGateway() *runtime.ServeMux {
	return runtime.NewServeMux(
		runtime.WithOutgoingHeaderMatcher(s.isHeaderAllowed),
		runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			header := request.Header.Get("Authorization")
			md := metadata.Pairs("auth", header)
			return md
		}),
		runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, writer http.ResponseWriter, request *http.Request, err error) {
			newError := runtime.HTTPStatusError{
				HTTPStatus: 400,
				Err:        err,
			}
			runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, writer, request, &newError)
		}),
	)
}

func (*Server) isHeaderAllowed(s string) (string, bool) {
	if _, isAllowed := allowedHeaders[s]; isAllowed {
		return strings.ToUpper(s), true
	}
	return s, false
}
