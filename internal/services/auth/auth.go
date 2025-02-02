package authservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Woland-prj/microtasks_sso/internal/domain/cerrors"
	"github.com/Woland-prj/microtasks_sso/internal/domain/dtos"
	"github.com/Woland-prj/microtasks_sso/internal/domain/entities"
	"github.com/Woland-prj/microtasks_sso/internal/lib/jwt"
	"github.com/Woland-prj/microtasks_sso/internal/lib/logger/sl"
	"golang.org/x/crypto/bcrypt"
)

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		user *entities.User,
	) (int64, error)
}

type UserProvider interface {
	GetUser(
		ctx context.Context,
		email string,
	) (*entities.User, error)
}

type AppProvider interface {
	GetApp(
		ctx context.Context,
		id int64,
	) (*entities.App, error)
}

type AuthService struct {
	log             *slog.Logger
	userSaver       UserSaver
	userProvider    UserProvider
	appProvider     AppProvider
	authTokenTTL    time.Duration
	refreshTokenTTL time.Duration
}

// New returns new AuthService instance
func New(
	log *slog.Logger,
	authTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
) *AuthService {
	return &AuthService{
		log:             log,
		userSaver:       userSaver,
		userProvider:    userProvider,
		appProvider:     appProvider,
		authTokenTTL:    authTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

// Register checks if user exists and if not exists, registers new user.
//
// If user exists, returns error.
// If user doesn't exist, creates user and saves to storage, returns uid.
func (a *AuthService) Register(
	ctx context.Context,
	dto dtos.RegisterDto,
) (int64, error) {
	const op = "authservice.Register"

	a.log.With(slog.String("op", op))
	a.log.Debug("registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		a.log.Error("failed to generate password hash", sl.Err(err))
		return 0, fmt.Errorf(
			"%s: %w",
			op,
			cerrors.NewCriticalInternalError("bcrypt.GenerateFromPassword", err),
		)
	}

	usr := &entities.User{
		Email:    dto.Email,
		PassHash: string(passHash),
	}

	uid, err := a.userSaver.SaveUser(ctx, usr)
	if err != nil {
		if errors.Is(err, &cerrors.AlreadyExistsError{}) {
			a.log.Warn("user exists", sl.Err(err))
			return 0, err
		}
		a.log.Error("failed to save user", sl.Err(err))
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	a.log.Debug("user registerd", slog.String("uid", fmt.Sprintf("%v", uid)))

	return uid, nil
}

// Login checks if user exists and if exists, returns pair of JWT tokens.
// If user doesn't exist, returns error.
func (a *AuthService) Login(
	ctx context.Context,
	dto dtos.LoginDto,
) (*entities.JwtTokenPair, error) {
	const op = "authservice.Login"

	a.log.With(slog.String("op", op))
	a.log.Debug("login user")

	usr, err := a.userProvider.GetUser(ctx, dto.Email)
	if err != nil {
		if errors.Is(err, &cerrors.NotFoundError{}) {
			a.log.Warn("user not found", sl.Err(err))
			return nil, cerrors.NewInvalidCredentialsError()
		}
		a.log.Error("failed to get user from storage", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usr.PassHash), []byte(dto.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			a.log.Warn("password mismatch", sl.Err(err))
			return nil, cerrors.NewInvalidCredentialsError()
		}
		a.log.Error("failed to compare password", sl.Err(err))
		return nil, fmt.Errorf(
			"%s: %w",
			op,
			cerrors.NewCriticalInternalError("bcrypt.CompareHashAndPassword", err),
		)
	}

	app, err := a.appProvider.GetApp(ctx, dto.AppId)
	if err != nil {
		if errors.Is(err, &cerrors.NotFoundError{}) {
			a.log.Warn("app not found", sl.Err(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		a.log.Error("failed to get app from storage", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	a.log.Debug("user logged in successfully", slog.String("uid", fmt.Sprintf("%v", usr.UID)))

	a.log.Debug("generating tokens")

	tokens, err := jwt.NewTokenPair(usr, app, a.authTokenTTL, a.refreshTokenTTL)
	if err != nil {
		a.log.Error("failed to generate tokens", sl.Err(err))
		return nil, fmt.Errorf(
			"%s: %w",
			op,
			cerrors.NewCriticalInternalError("jwt.NewTokenPair", err),
		)
	}

	a.log.Debug(
		"tokens generated",
		slog.String("auth", tokens.AuthToken),
		slog.String("refresh", tokens.RefreshToken),
	)

	return tokens, nil
}
