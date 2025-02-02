package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Woland-prj/microtasks_sso/internal/app"
	"github.com/Woland-prj/microtasks_sso/internal/config"
	"github.com/Woland-prj/microtasks_sso/internal/lib/logger/handlers/slogpretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// Load config
	conf := config.MustLoad()

	// Setup logger
	log := mustSetupLogger(conf.Env)
	log.Info("Starting app...")
	log.Debug("Logger init")

	// Run app
	application := app.New(
		log,
		conf.GRPC.Port,
		conf.StoragePath,
		conf.TokenTTL.Auth,
		conf.TokenTTL.Refresh,
	)

	go application.GRPCServ.MustRun()

	// Graceful shutdown
	chanStop := make(chan os.Signal, 1)
	signal.Notify(chanStop, syscall.SIGTERM, syscall.SIGINT)

	stopSignal := <-chanStop
	log.Info("Trying stop application", slog.String("signal", stopSignal.String()))

	application.GRPCServ.Stop()
	log.Info("Application stopped")
}

func mustSetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		panic("Unknown env: " + env)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
