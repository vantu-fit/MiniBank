package gapi

import (
	"context"

	db "github.com/vantu-fit/master-go-be/db/sqlc"
	"github.com/vantu-fit/master-go-be/pb"
	"github.com/vantu-fit/master-go-be/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	violations := validateVerifyEmailRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	result , err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailId:    req.EmailId,
		SecretCode: req.SecretCode,
	})
	if err != nil {
		return nil , status.Errorf(codes.Internal , "falies to verify email")
	}
	response := &pb.VerifyEmailResponse{
		IsVerified: result.User.IsEmailVerified,
	}
	return response, nil

}

func validateVerifyEmailRequest(req *pb.VerifyEmailRequest) (violatios []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmailId(req.EmailId); err != nil {
		violatios = append(violatios, fileViolation("email_id", err))
	}
	if err := val.ValidateSecretCode(req.SecretCode); err != nil {
		violatios = append(violatios, fileViolation("secret_code", err))
	}

	return violatios

}
