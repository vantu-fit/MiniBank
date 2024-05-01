package gapi

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	db "github.com/vantu-fit/master-go-be/db/sqlc"
	"github.com/vantu-fit/master-go-be/pb"
	"github.com/vantu-fit/master-go-be/utils"
	"github.com/vantu-fit/master-go-be/val"
	"github.com/vantu-fit/master-go-be/worker"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// TODO : understand this code block
	var response *pb.CreateUserResponse
	violations := validateCreteUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	hassPassword, err := utils.HashedPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hassPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}
	user, err := server.store.CreateUserTx(ctx, db.CreateUserTxParams{
		CreateUserParams: arg,
		AfterCreate: func(user db.User) error {
			// create task send verify email
			payload := &worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}

			opts := []asynq.Option{
				asynq.MaxRetry(10),                // thu lai 10 lan
				asynq.ProcessIn(10 * time.Second), // do tre 10s
				asynq.Queue(worker.QueueCritical), // hang doi uu tien
			}

			err = server.taskDitributor.DistributeTaskSendVerifyEmail(ctx, payload, opts...)
			if err != nil {
				return status.Errorf(codes.Internal, "falied to dostribute task to send verify email: %s", err)
			}
			return nil
		},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create user: %s", err)
	}

	response = &pb.CreateUserResponse{
		User: &pb.User{
			Username:          user.User.Username,
			FullName:          user.User.FullName,
			Email:             user.User.Email,
			PasswordChangedAt: timestamppb.New(user.User.PasswordChangedAt),
			CreatedAt:         timestamppb.New(user.User.CreatedAt),
		},
	}
	fmt.Println("done create user >>>")

	return response, nil

}

func validateCreteUserRequest(req *pb.CreateUserRequest) (violatios []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violatios = append(violatios, fileViolation("username", err))
	}
	if err := val.ValidateFullname(req.GetFullName()); err != nil {
		violatios = append(violatios, fileViolation("full_name", err))
	}
	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violatios = append(violatios, fileViolation("email", err))
	}
	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violatios = append(violatios, fileViolation("password", err))
	}

	return violatios
}
