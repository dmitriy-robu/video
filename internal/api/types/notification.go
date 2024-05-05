package types

import (
	"go-fitness/internal/api/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Notification struct {
	ID        primitive.ObjectID      `bson:"_id,omitempty"`
	UUID      string                  `bson:"uuid"`
	Name      string                  `bson:"name"`
	Body      string                  `bson:"body"`
	Status    enum.NotificationStatus `bson:"type"`
	IsRead    bool                    `bson:"is_read"`
	CreatedAt time.Time               `bson:"created_at"`
	UpdatedAt time.Time               `bson:"updated_at"`
}
