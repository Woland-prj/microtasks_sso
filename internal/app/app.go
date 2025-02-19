package app

import (
	"log/slog"

	grpc_app "github.com/Woland-prj/microtasks_sso/internal/app/grpc"
	http_app "github.com/Woland-prj/microtasks_sso/internal/app/http"
	"github.com/Woland-prj/microtasks_sso/internal/config"
	"github.com/Woland-prj/microtasks_sso/internal/lib/logger/sl"
	"github.com/Woland-prj/microtasks_sso/internal/services"
	"github.com/Woland-prj/microtasks_sso/internal/storage/sqlite"
	"github.com/go-playground/validator/v10"
)

type App struct {
	log      *slog.Logger
	srvs     *services.Services
	GRPCServ *grpc_app.App
	HTTPServ *http_app.App
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	storage := mustCreateSqliteStorage(log, cfg.StoragePath)
	services := services.New(
		log, 
		storage, 
		cfg.TokenTTL.Auth, 
		cfg.TokenTTL.Refresh,
	)
	validate := validator.New(validator.WithRequiredStructEnabled())
	grpcapp := grpc_app.New(log, cfg.GRPC.Port, services, validate)
	httpapp := http_app.New(
		log, 
		cfg.HTTP.Port, 
		cfg.HTTP.Timeout, 
		cfg.HTTP.IdleTimeout, 
		cfg.HTTP.StopTimeout, 
		services, 
		validate,
	)

	return &App{
		log:      log,
		srvs:     services,
		GRPCServ: grpcapp,
		HTTPServ: httpapp,
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
