package gapi

import (
	"context"
	"database/sql"
	"time"

	db "github.com/longln/simplebank/db/sqlc"
	"github.com/longln/simplebank/pb"
	"github.com/longln/simplebank/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if authPayload.Username != request.GetUserName() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's information")
	}

	arg := db.UpdateUserParams{
		Username: request.GetUserName(),
		FullName: sql.NullString{String: request.GetFullName(), Valid: request.FullName != nil},
		Email: sql.NullString{String: request.GetEmail(), Valid: request.Email != nil},
	}

	if request.Password != nil {
		hashedpassword, err := utils.HashPassword(request.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}
		arg.HashedPassword = sql.NullString{String: hashedpassword, Valid: true}

		arg.PasswordChangedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to Update user: %s", err)
	}
	response := &pb.UpdateUserResponse{
		User: convertUser(user),
	}
	return response, nil
}


