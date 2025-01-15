package app

import (
	"log/slog"
	"time"

	grpc_app "github.com/Woland-prj/microtasks_sso/internal/app/grpc"
)

type App struct {
	log      *slog.Logger
	GRPCServ *grpc_app.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	grpcapp := grpc_app.New(log, grpcPort)

	return &App{
		log:      log,
		GRPCServ: grpcapp,
	}
}
