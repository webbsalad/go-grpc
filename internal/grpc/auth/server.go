package auth

import (
	"context"

	appv1 "github.com/webbsalad/test-protos/gen/go/app"
	"google.golang.org/grpc"
)

type serverAPI struct {
	appv1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	appv1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *appv1.LoginRequest,
) (*appv1.LoginResponse, error) {
	panic("1")
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
