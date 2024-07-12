package auth

import (
	"context"

	appv1 "github.com/webbsalad/test-protos/gen/go/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)

	RegisterNewUser(ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)

	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	appv1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	appv1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *appv1.LoginRequest,
) (*appv1.LoginResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "not email")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "not password")
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.Password, int(req.GetAppId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "iternal error")

	}

	return &appv1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *appv1.RegisterRequest,
) (*appv1.RegisterResponse, error) {
	panic("1")
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	req *appv1.IsAdminRequest,
) (*appv1.IsAdminResponse, error) {
	panic("1")
}
