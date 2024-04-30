package repository

import (
	"go.uber.org/fx"
)

func NewRepository() fx.Option {
	return fx.Module(
		"repository",
		fx.Provide(
			fx.Annotate(
				NewVideoRepository,
				fx.As(new(VideoRepositoryInterface)),
			),

			fx.Annotate(
				NewNotificationRepository,
				fx.As(new(NotificationRepoInterface)),
			),

			fx.Annotate(
				NewUserRepository,
				fx.As(new(UserRepositoryInterface)),
			),
		),
	)
}
