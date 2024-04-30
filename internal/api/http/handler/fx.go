package handler

import (
	"go.uber.org/fx"
)

type Handlers struct {
	Video *VideoHandler
}

func NewHandlers(
	Video *VideoHandler,
) *Handlers {
	return &Handlers{
		Video: Video,
	}
}

func NewHandler() fx.Option {
	return fx.Module(
		"handler",
		fx.Options(),
		fx.Provide(
			NewVideoHandler,
			NewHandlers,
		),
	)
}
