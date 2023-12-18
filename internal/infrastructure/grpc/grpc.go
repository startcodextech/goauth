package grpc

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	server "github.com/startcodextech/goauth/internal/infrastructure/http"
	"github.com/startcodextech/goauth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"net"
	"net/http"
	"strings"
)

const (
	grpcNetwork = "tcp"
	grpcAddress = "0.0.0.0:8080"
)

var (
	allowedHeaders = map[string]struct{}{
		"x-request-id": {},
	}
)

func Start(ctx context.Context, commandBus *cqrs.CommandBus, eventSubscriber message.Subscriber, logger watermill.LoggerAdapter) {

	rpcServer := grpc.NewServer()
	listenRpc, err := net.Listen(grpcNetwork, grpcAddress)
	if err != nil {
		logger.Error("", err, nil)
		panic(err)
	}

	go func(server *grpc.Server, listen net.Listener, commandBus *cqrs.CommandBus, eventSubscriber message.Subscriber, logger watermill.LoggerAdapter) {

		accountService := NewAccountService(commandBus, eventSubscriber, logger)
		proto.RegisterAccountServiceServer(server, accountService)

		err := server.Serve(listen)
		if err != nil {
			logger.Error("", err, nil)
			panic(err)
		}
	}(rpcServer, listenRpc, commandBus, eventSubscriber, logger)

	mux := runtime.NewServeMux(
		runtime.WithOutgoingHeaderMatcher(isHeaderAllowed),
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

	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err = proto.RegisterAccountServiceHandlerFromEndpoint(ctx, mux, grpcAddress, dialOptions)
	if err != nil {
		logger.Error("", err, nil)
		panic(err)
	}

	serverHttp := server.New()
	serverHttp.Group("api/v1/*", adaptor.HTTPHandler(mux))
	go server.Start(serverHttp, logger)

}

func isHeaderAllowed(s string) (string, bool) {
	if _, isAllowed := allowedHeaders[s]; isAllowed {
		return strings.ToUpper(s), true
	}
	return s, false
}
