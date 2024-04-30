package service

import (
	"go.uber.org/fx"
)

type Service interface {
}

func NewService() fx.Option {
	return fx.Module(
		"service",
		fx.Provide(
			NewVideoTranscodeTask,

			fx.Annotate(
				NewVideoService,
				fx.As(new(VideoServiceInterface)),
				fx.As(new(UploadAndTranscodeQueueInterface)),
			),

			fx.Annotate(
				NewNotificationService,
				fx.As(new(NotificationServiceInterface)),
			),
		),
	)
}
