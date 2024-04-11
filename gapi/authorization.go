package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/vantu-fit/master-go-be/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader     = "authorization"
	authorizationTypeBearer = "bearer"
)

func (server *Server) authorizationUser(ctx context.Context , accessibleRole []string) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationTypeBearer {
		return nil, fmt.Errorf("unsupport authorization type: %s", authType)
	}

	accessToken := fields[1]
	payload, err := server.maker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token")
	}

	if !hassPermission(payload.Role , accessibleRole) {
		return nil , fmt.Errorf("permission denied")
	}

	return payload, nil

}

func hassPermission(userRole string , accessibleRole []string) bool {
	for _ , role := range accessibleRole {
		if userRole == role {
			return true
		}
	}
	return false
}
