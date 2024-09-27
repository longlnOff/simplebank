package gapi

import (
	"context"
	"database/sql"

	db "github.com/longln/simplebank/db/sqlc"
	"github.com/longln/simplebank/pb"
	"github.com/longln/simplebank/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)


func (server *Server) LoginUser(ctx context.Context, request *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, request.UserName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find user")

	}
	err = utils.CheckHashPassword(request.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password")
	}
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.UserName, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token")
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.UserName, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token")
	}
	// create session
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserName:     user.UserName,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,

	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session")
	}


	rsp := &pb.LoginUserResponse{
		User: convertUser(user),
		SessionId: session.ID.String(),
		AccessToken: accessToken,
		AccessTokenExpiresAt: timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken: refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}
	
	return rsp, nil
}