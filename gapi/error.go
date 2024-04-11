package gapi

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func fileViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

func invalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
	badRequest := &errdetails.BadRequest{FieldViolations: violations}
	statusIvalid := status.New(codes.InvalidArgument, "invalid parameter")

	statusDetails, err := statusIvalid.WithDetails(badRequest)
	if err != nil {
		return statusIvalid.Err()
	}

	return statusDetails.Err()
}
