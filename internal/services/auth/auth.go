package auth

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"sso/internal/domain/models"
	"sso/internal/lib/jwt"
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
	ErrInvalidAppId       = errors.New("invalid app id")
	ErrUserExist          = errors.New("user exist")
)

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, id int32) (models.App, error)
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

func (auth *Auth) Login(ctx context.Context, email, password string, appId int32) (string, error) {
	const op = "auth.Login"

	log := auth.log.With(slog.String("op", op), slog.String("username", email))

	log.Info("logging user")

	user, err := auth.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			auth.log.Warn("user not found")

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		auth.log.Warn("failed to get user", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		auth.log.Info("invalid credentials", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := auth.appProvider.App(ctx, appId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successfully")

	token, err := jwt.NewToken(user, app, auth.tokenTTL)
	if err != nil {
		auth.log.Error("failed to generate token", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
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
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", sl.Err(err))

			return 0, fmt.Errorf("%s: %w", op, ErrUserExist)
		}

		log.Error("failed to save new user", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// IsAdmin checks if user is admin.
func (auth *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "Auth.IsAdmin"

	log := auth.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)

	log.Info("checking if user is admin")

	isAdmin, err := auth.usrProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("user not found")

			return false, fmt.Errorf("%s: %w", op, ErrInvalidAppId)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
