package v1

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes"
	v1 "github.com/namtx/go-grpc-rest-microservice-tutorial/pkg/api/v1"
)

const (
	apiVersion = "v1"
)

type todoServiceServer struct {
	db *sql.DB
}

func NewTodoServiceServer(db *sql.DB) v1.TodoServiceServer {
	return &todoServiceServer{db: db}
}

func (s *todoServiceServer) checkAPI(api string) error {
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented, "unsupported API version: service implements API version %s, but asked for %s", apiVersion, api)
		}
	}

	return nil
}

// connnect method returns SQL database connetion from the pool
func (s *todoServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database: "+err.Error())
	}

	return c, nil
}

// Create Todo task
func (s *todoServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	reminder, err := ptypes.Timestamp(req.Todo.Reminder)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "reminder has invalid format: "+err.Error())
	}

	// insert Todo entity data
	res, err := c.ExecContext(ctx, "INSERT INTO todos(`Title`, `Description`, `Reminder`) VALUES (?, ?, ?)", req.Todo.Title, req.Todo.Description, reminder)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into Todo: "+err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve last insert id: "+err.Error())
	}

	return &v1.CreateResponse{
		Api: apiVersion,
		Id:  id,
	}, nil
}
