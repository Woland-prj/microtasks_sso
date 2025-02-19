package http_app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	authhttp "github.com/Woland-prj/microtasks_sso/internal/http/auth"
	mvLogger "github.com/Woland-prj/microtasks_sso/internal/http/middleware/logger"
	"github.com/Woland-prj/microtasks_sso/internal/lib/logger/sl"
	"github.com/Woland-prj/microtasks_sso/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

type App struct{
	log *slog.Logger
	server *http.Server
	port int
	stopTimeout time.Duration
}

func New(
	log *slog.Logger,
	port int,
	timeout time.Duration,
	idleTimeout time.Duration,
	stopTimeout time.Duration,
	services *services.Services,
	validate *validator.Validate,
) *App{
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(mvLogger.New(log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	authhttp.Register(r, services.Auth, validate)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: r,
		ReadTimeout: timeout,
		WriteTimeout: timeout,
		IdleTimeout: idleTimeout,
	}
	return &App{log: log, server: srv, port: port}
}

func (a *App) MustRun(){
	if err := a.run(); err != nil {
		panic(err)
	}
}

func (a *App) run() error {
	op := "http_app.Run"
	a.log = a.log.With(slog.String("op", op))
	a.log.Info("Starting HTTP server", slog.Int("port", a.port))

	a.log.Info("HTTP server is runing", slog.String("addr", a.server.Addr))
	
	err := a.server.ListenAndServe()
	if err != nil  && err != http.ErrServerClosed {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "grpc_app.Stop"

	a.log.With(slog.String("op", op)).
		Info("Stopping HTTP server", slog.Int("port", a.port))

	ctx, cancel := context.WithTimeout(context.Background(), a.stopTimeout)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Error("failed to shutdown HTTP server", sl.Err(err))
	}
}