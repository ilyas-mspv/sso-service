package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso-service/internal/app"
	"sso-service/internal/config"
	"syscall"
)

var (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting app")

	application := app.New(log, cfg.GRPC.Port, cfg.DatabaseURL, cfg.TokenTTl)
	go application.GRPCSrv.MustRun()

	gracefulShutdown(log, application)
}

func gracefulShutdown(log *slog.Logger, application *app.App) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	signal := <-stop
	application.GRPCSrv.Stop()
	log.Info("Application stopped", slog.String("signal", signal.String()))
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:

		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
