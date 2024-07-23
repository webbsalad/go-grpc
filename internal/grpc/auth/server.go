package auth

import (
	"context"
	"errors"

	"github.com/webbsalad/go-grpc/internal/services/auth"
	"github.com/webbsalad/go-grpc/internal/storage"
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
	in *appv1.LoginRequest,
) (*appv1.LoginResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "почта обязательна")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "пароль бязателен")
	}

	if in.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "app_id обязателен")
	}

	token, err := s.auth.Login(ctx, in.GetEmail(), in.GetPassword(), int(in.GetAppId()))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "неверная почта или пароль")
		}

		return nil, status.Error(codes.InvalidArgument, "ошибка входа")
	}

	return &appv1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	in *appv1.RegisterRequest,
) (*appv1.RegisterResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "почта обязательна")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "пароль обязателен")
	}

	uid, err := s.auth.RegisterNewUser(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "пользователь уже существует")
		}

		return nil, status.Error(codes.Internal, "ошибка регистрации польщзователя")
	}

	return &appv1.RegisterResponse{UserId: uid}, nil
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	in *appv1.IsAdminRequest,
) (*appv1.IsAdminResponse, error) {
	if in.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id обязателен")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, in.GetUserId())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "пользователь не найден")
		}

		return nil, status.Error(codes.Internal, "ошибка проверки статуса админа")
	}

	return &appv1.IsAdminResponse{IsAdmin: isAdmin}, nil
}
