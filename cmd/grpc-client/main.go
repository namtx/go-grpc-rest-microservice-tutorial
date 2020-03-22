package main

import (
	"context"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	v1 "github.com/namtx/go-grpc-rest-microservice-tutorial/pkg/api/v1"
	"google.golang.org/grpc"
)

const (
	apiVersion = "v1"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	c := v1.NewTodoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := time.Now().In(time.UTC)
	reminder, _ := ptypes.TimestampProto(t)
	pfx := t.Format(time.RFC3339Nano)

	request := v1.CreateRequest{
		Api: apiVersion,
		Todo: &v1.Todo{
			Title:       "title (" + pfx + ")",
			Description: "description (" + pfx + ")",
			Reminder:    reminder,
		},
	}

	response, err := c.Create(ctx, &request)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}

	log.Printf("Create result: <+%v>\n\n", response)
}
