package app

import (
	"log/slog"
	grpcapp "sso-service/internal/app/grpc"
	"sso-service/internal/services/auth"
	"sso-service/internal/storage/pgsql"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	databaseURL string,
	tokenTTL time.Duration,
) *App {
	storage, err := pgsql.New(databaseURL)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)
	return &App{
		GRPCSrv: grpcApp,
	}
}
