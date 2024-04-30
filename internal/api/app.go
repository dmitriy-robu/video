package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/patrickmn/go-cache"
	"github.com/pusher/pusher-http-go/v5"
	"go-fitness/external/config"
	"go-fitness/external/db"
	"go-fitness/internal/api/event"
	"go-fitness/internal/api/http/handler"
	"go-fitness/internal/api/http/middleware"
	"go-fitness/internal/api/repository"
	"go-fitness/internal/api/service"
	"go.uber.org/fx"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Options(
			repository.NewRepository(),
			service.NewService(),
			handler.NewHandler(),
			middleware.NewMiddleware(),
			db.NewDataBase(),
		),
		fx.Provide(
			config.NewConfig,
			fx.Annotate(
				event.NewPusherEvent,
				fx.As(new(event.WSInterface)),
			),
			validator.New,
			NewCache,
			NewLogger,
			NewRouter,
			NewServer,
			NewPusher,
		),
		fx.Invoke(RunServer),
	)
}

func NewPusher(cfg *config.Config) pusher.Client {
	return pusher.Client{
		AppID:   cfg.WSServer.AppID,
		Key:     cfg.WSServer.Key,
		Secret:  cfg.WSServer.Secret,
		Cluster: cfg.WSServer.Cluster,
		Host:    cfg.WSServer.Host + ":" + cfg.WSServer.Port,
		Secure:  cfg.WSServer.Secure,
	}
}

func NewCache() *cache.Cache {
	return cache.New(cache.NoExpiration, cache.NoExpiration)
}
