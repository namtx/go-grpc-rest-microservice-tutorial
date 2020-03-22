package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/namtx/go-grpc-rest-microservice-tutorial/pkg/protocol/grpc"
	v1 "github.com/namtx/go-grpc-rest-microservice-tutorial/pkg/service/v1"
)

type Todo struct {
	gorm.Model
	Title       string
	Description string
	Reminder    time.Time
}

func RunServer() error {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", "todo_service.db")

	if err != nil {
		return fmt.Errorf("Failed to open database %v", err)
	}
	defer db.Close()

	v1API := v1.NewTodoServiceServer(db)

	return grpc.RunServer(ctx, v1API, "50051")
}
