package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/webbsalad/go-grpc/internal/domain/models"
	"github.com/webbsalad/go-grpc/internal/lib/jwt"
	"github.com/webbsalad/go-grpc/internal/lib/logger/sl"
	"github.com/webbsalad/go-grpc/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		log:         log,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appID int,
) (string, error) {
	const op = "Auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", email),
	)

	log.Info("попытка залогинеть юзера")

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("пользователь не найден", sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, "ErrInvalidCredentials")
		}

		a.log.Error("ошибка при получении юзера", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("неверные данные для входа", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, "ErrInvalidCredentials")
	}

	app, err := a.appProvider.App(ctx, appID)
	fmt.Println(err)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("пользователь залогинелся")

	token, err := jwt.NerToken(user, app, a.tokenTTL)

	if err != nil {
		a.log.Error("ошибка генерации токена", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil

}

func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	pass string,
) (int64, error) {
	const op = "auth.RegisterNewUser"
	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	log.Info("регистрация пользователя")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	if err != nil {
		log.Error("проблемы с генерацией пароля", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)

	if err != nil {
		log.Error("ошибка сохранения юзера", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (a *Auth) IsAdmin(
	ctx context.Context,
	userID int64,
) (bool, error) {
	const op = "Auth.IsAdmin"

	log.Info("проверка если юзер админ")

	IsAdmin, err := a.usrProvider.IsAdmin(ctx, userID)

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("проверено сли пользователь админ", slog.Bool("is_admin: ", IsAdmin))

	return IsAdmin, nil

}
