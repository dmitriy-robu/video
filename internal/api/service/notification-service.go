package service

import (
	"context"
	"errors"
	"go-fitness/external/logger/sl"
	"go-fitness/internal/api/event"
	"go-fitness/internal/api/repository"
	"go-fitness/internal/api/types"
	"log/slog"
)

type NotificationService struct {
	log              *slog.Logger
	event            event.WSInterface
	notificationRepo repository.NotificationRepoInterface
}

type NotificationServiceInterface interface {
	ProcessNotification(ctx context.Context, notification types.Notification) error
}

func NewNotificationService(
	log *slog.Logger,
	event event.WSInterface,
	notificationRepo repository.NotificationRepoInterface,
) *NotificationService {
	return &NotificationService{
		log:              log,
		event:            event,
		notificationRepo: notificationRepo,
	}
}

func (s *NotificationService) ProcessNotification(ctx context.Context, notification types.Notification) error {
	const op = "NotificationService.ProcessNotification"

	log := s.log.With(
		sl.String("op", op),
	)

	log.Info("processing notification")

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		log.Error("can't create notification", sl.Err(err))
		return errors.New("can't create notification")
	}

	go func() {
		if err := s.sendNotificationByEmail(notification); err != nil {
			log.Error("can't send notification by email", sl.Err(err))
		}
	}()

	go func() {
		if err := s.sendNotificationByEvent(notification); err != nil {
			log.Error("can't send notification by event", sl.Err(err))
		}
	}()

	return nil
}

func (s *NotificationService) sendNotificationByEmail(notification types.Notification) error {
	const op = "NotificationService.sendNotificationByEmail"

	log := s.log.With(
		sl.String("op", op),
	)

	log.Info("sending notification by email")

	return nil
}

func (s *NotificationService) sendNotificationByEvent(notification types.Notification) error {
	const op = "NotificationService.sendNotificationByEvent"

	log := s.log.With(
		sl.String("op", op),
	)

	log.Info("sending notification by event")

	if err := s.event.TriggerEvent(event.NewNotificationEvent(notification)); err != nil {
		log.Error("failed to trigger event", sl.Err(err))
		return errors.New("failed to trigger event")
	}

	return nil
}
