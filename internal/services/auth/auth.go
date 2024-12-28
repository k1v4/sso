package auth

import (
	"context"
	"errors"
	"log/slog"
	"sso/internal/domain/models"
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
	SaveUser(ctx context.Context, email string, passHash []byte) (int64, error)
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
	panic("implement me")
}

func (auth *Auth) RegisterNewUser(ctx context.Context, email, password string) (*models.User, error) {
	panic("implement me")
}

func (auth *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	panic("implement me")
}
