package grpc

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/startcodextech/goauth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
)

const (
	grpcNetwork = "tcp"
	grpcAddress = "0.0.0.0:9090"
	httpAddress = "0.0.0.0:8000"
)

func Start(ctx context.Context, commandBus *cqrs.CommandBus, eventSubscriber message.Subscriber, logger watermill.LoggerAdapter) {

	listen, err := net.Listen(grpcNetwork, grpcAddress)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	rpcServer := grpc.NewServer()

	accountService := NewAccountService(commandBus, eventSubscriber, logger)
	proto.RegisterAccountServiceServer(rpcServer, accountService)

	go func() {
		err := rpcServer.Serve(listen)
		if err != nil {
			logger.Error("", err, nil)
		}
	}()

	conn, err := grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {

	}
	defer func(conn *grpc.ClientConn, logger watermill.LoggerAdapter) {
		err := conn.Close()
		if err != nil {
			logger.Error("", err, nil)
		}
	}(conn, logger)

	mux := runtime.NewServeMux()

	err = proto.RegisterAccountServiceHandler(ctx, mux, conn)
	if err != nil {
		if err != nil {
			logger.Error("", err, nil)
		}
	}

	gwServer := &http.Server{
		Addr:    httpAddress,
		Handler: mux,
	}

	logger.Info("Serving gRPC-Gateway on connection", nil)

	go func(server *http.Server) {
		log.Fatalln(server.ListenAndServe())
	}(gwServer)
}
