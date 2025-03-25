package gapi

import (
	"context"
	"database/sql"

	db "github.com/phongnd2802/simple-bank/db/sqlc"
	"github.com/phongnd2802/simple-bank/pb"
	"github.com/phongnd2802/simple-bank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "username %s not found", req.GetUsername())
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %s", err)
	}

	err = util.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "username or password is incorrect")
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.Token.AccessTokenDuration)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create access token")
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.Token.RefreshTokenDuration)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create refresh token")
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		RefreshToken: refreshToken,
		Username:     user.Username,
		UserAgent:    "",
		ClientIp:     "",
		ExpiresAt:    refreshPayload.ExpiresAt.Time,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %s", err)
	}

	rsp := &pb.LoginUserResponse{
		SessionId:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiresAt.Time),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiresAt.Time),
		User: convertUser(user),
	}

	return rsp, nil
}
