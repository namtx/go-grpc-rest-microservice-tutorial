package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	v1 "github.com/namtx/go-grpc-rest-microservice-tutorial/pkg/api/v1"
)

func RunServer(ctx context.Context, v1API v1.TodoServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()

	v1.RegisterTodoServiceServer(server, v1API)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			log.Println("shuting down gRPC server...")
			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	log.Println("starting gRPC server...")

	return server.Serve(listen)
}
