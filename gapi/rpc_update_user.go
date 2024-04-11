package gapi

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/vantu-fit/master-go-be/db/sqlc"
	"github.com/vantu-fit/master-go-be/pb"
	"github.com/vantu-fit/master-go-be/utils"
	"github.com/vantu-fit/master-go-be/val"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// TODO : understand this code block
	payload, err := server.authorizationUser(ctx, []string{utils.BankerRole, utils.DepositorRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	arg := db.UpdateUserParams{
		Username: req.Username,
		FullName: pgtype.Text{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
		Email: pgtype.Text{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
		HashedPassword: pgtype.Text{
			String: req.GetPassword(),
			Valid:  req.Password != nil,
		},
	}

	if req.Username != payload.Username && payload.Role != utils.BankerRole { 
		return nil, status.Errorf(codes.InvalidArgument, "invalid username %s", err)
	}

	if req.Password != nil {
		hassPassword, err := utils.HashedPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "fail to hash password %s", err)
		}
		arg.HashedPassword = pgtype.Text{
			String: hassPassword,
			Valid:  true,
		}

		arg.PasswordChangedAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: req.Password != nil,
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "cannot upadte user: %s", err)
	}

	response := pb.UpdateUserResponse{
		User: &pb.User{
			Username:          user.Username,
			FullName:          user.FullName,
			Email:             user.Email,
			PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
			CreatedAt:         timestamppb.New(user.CreatedAt),
		},
	}

	return &response, nil

}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violatios []*errdetails.BadRequest_FieldViolation) {

	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violatios = append(violatios, fileViolation("username", err))
	}

	if req.FullName != nil {
		if err := val.ValidateFullname(req.GetFullName()); err != nil {
			violatios = append(violatios, fileViolation("full_name", err))
		}
	}

	if req.Email != nil {
		if err := val.ValidateEmail(req.GetEmail()); err != nil {
			violatios = append(violatios, fileViolation("email", err))
		}
	}
	if req.Password != nil {
		if err := val.ValidatePassword(req.GetPassword()); err != nil {
			violatios = append(violatios, fileViolation("password", err))
		}
	}

	return violatios

}

func unauthenticatedError(err error) error {
	return status.Errorf(codes.Unauthenticated, "unauthenticated error: %s", err)
}
