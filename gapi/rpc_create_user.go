package gapi

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	db "github.com/longln/simplebank/db/sqlc"
	"github.com/longln/simplebank/pb"
	"github.com/longln/simplebank/utils"
	"github.com/longln/simplebank/worker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedpassword, err := utils.HashPassword(request.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}
	
	arg := db.CreateUserParams{
		UserName: request.GetUserName(),
		HashedPassword: hashedpassword,
		FullName: request.GetFullName(),
		Email: request.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "user name already exist: %s", err)

			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}


	// TODO use DB transaction because it's not atomic
	// send email verification
	payload := worker.PayloadSendVerifyEmail{
		Username: user.UserName,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(5),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.CriticalQueue),
	}
	err = server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, payload, opts...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to distribute send verify email: %s", err)
	}
	response := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	
	return response, nil
}


