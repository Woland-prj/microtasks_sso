package app

import (
	"log/slog"
	"time"

	grpc_app "github.com/Woland-prj/microtasks_sso/internal/app/grpc"
	"github.com/Woland-prj/microtasks_sso/internal/lib/logger/sl"
	"github.com/Woland-prj/microtasks_sso/internal/services"
	"github.com/Woland-prj/microtasks_sso/internal/storage/sqlite"
	"github.com/go-playground/validator/v10"
)

type App struct {
	log      *slog.Logger
	srvs     *services.Services
	GRPCServ *grpc_app.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	authTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *App {
	storage := mustCreateSqliteStorage(log, storagePath)
	services := services.New(log, storage, authTokenTTL, refreshTokenTTL)
	validate := validator.New(validator.WithRequiredStructEnabled())
	grpcapp := grpc_app.New(log, grpcPort, services, validate)

	return &App{
		log:      log,
		srvs:     services,
		GRPCServ: grpcapp,
	}
}

func mustCreateSqliteStorage(log *slog.Logger, storagePath string) *sqlite.Storage {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		log.Error("Error creating sqlite storage", sl.Err(err))
		panic("Error creating sqlite storage, see logs")
	}

	return storage
}
