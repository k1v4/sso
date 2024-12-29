package auth

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"sso/internal/domain/models"
	"sso/internal/lib/logger/sl"
	"sso/internal/storage"
	"time"
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

func New(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		log:         log,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

func (auth *Auth) Login(ctx context.Context, email, password string, appId int32) (*models.User, error) {
	const op = "auth.Login"

	log := auth.log.With(slog.String("op", op), slog.String("username", email))

	log.Info("logging user")

	user, err := auth.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			auth.log.Warn("user not found")

			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		auth.log.Warn("failed to get user", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		auth.log.Info("invalid credentials", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := auth.appProvider.App(ctx, appId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	auth.log.Info("user logged in successfully")
}

func (auth *Auth) RegisterNewUser(ctx context.Context, email, password string) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := auth.log.With(slog.String("op", op), slog.String("email", email))

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := auth.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save new user", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (auth *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	panic("implement me")
}
