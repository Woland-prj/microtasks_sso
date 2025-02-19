package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/Woland-prj/microtasks_sso/internal/domain/entities"
	authservice "github.com/Woland-prj/microtasks_sso/internal/services/auth"
)

type Services struct {
	Auth *authservice.AuthService
}

type Storage interface {
	SaveUser(
		ctx context.Context,
		user *entities.User,
	) (int64, error)

	GetUserByEmail(
		ctx context.Context,
		email string,
	) (*entities.User, error)

	GetUserById(
		ctx context.Context,
		uid int64,
	) (*entities.User, error)

	GetApp(
		ctx context.Context,
		id int64,
	) (*entities.App, error)
}

func New(
	log *slog.Logger,
	storage Storage,
	authTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *Services {
	return &Services{
		Auth: authservice.New(
			log,
			authTokenTTL,
			refreshTokenTTL,
			storage,
			storage,
			storage,
		),
	}
}
