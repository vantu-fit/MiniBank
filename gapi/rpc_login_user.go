package gapi

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	db "github.com/vantu-fit/master-go-be/db/sqlc"
	"github.com/vantu-fit/master-go-be/pb"
	"github.com/vantu-fit/master-go-be/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "User is no arealdy exist: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "Cannot get user: %s", err)
	}
	// check password with user password in database
	err = utils.CheckPassword(user.HashedPassword, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot get user: %s", err)

	}
	// parser duration in config file
	duration, err := time.ParseDuration(server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot parser duration: %s", err)

	}
	// create acctoken
	accessToken, accessPayload, err := server.maker.CreateToken(user.Username,user.Role , duration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot create access token: %s", err)

	}

	refresh_token_duration, err := time.ParseDuration(server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot parser refresh duration: %s", err)
	}

	refreshToken, refreshPayload, err := server.maker.CreateToken(user.Username,user.Role,  refresh_token_duration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot create refresh token: %s", err)

	}
	// extract metadata
	mtdt := server.extractMetadata(ctx)

	argSession := db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	}
	// create session
	session, err := server.store.CreateSession(ctx, argSession)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot parser duration: %s", err)

	}

	// create response
	response := pb.LoginUserResponse{
		AccessToken:           accessToken,
		AccessTokenExpriedAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpriedAt: timestamppb.New(refreshPayload.ExpiredAt),
		SessionId:             session.ID.String(),
		User: &pb.User{
			Username:          user.Username,
			FullName:          user.FullName,
			Email:             user.Email,
			PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
			CreatedAt:         (timestamppb.New(user.PasswordChangedAt)),
		},
	}

	return &response, nil
}
