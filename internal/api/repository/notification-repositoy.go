package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go-fitness/external/db"
	"go-fitness/internal/api/types"
	"time"
)

const NotificationCollection = "notifications"

type NotificationRepository struct {
	mdb db.MongoDBInterface
}

type NotificationRepoInterface interface {
	Create(ctx context.Context, notification types.Notification) error
}

func NewNotificationRepository(
	mdb db.MongoDBInterface,
) *NotificationRepository {
	return &NotificationRepository{
		mdb: mdb,
	}
}

func (r *NotificationRepository) Create(ctx context.Context, notification types.Notification) error {
	const op = "NotificationRepository.Create"

	now := time.Now()

	notification.UUID = uuid.New().String()
	notification.CreatedAt = now
	notification.UpdatedAt = now

	_, err := r.mdb.InsertOne(ctx, NotificationCollection, notification)
	if err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}

	return nil
}
