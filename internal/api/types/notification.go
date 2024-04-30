package types

import (
	"go-fitness/internal/api/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Notification struct {
	ID        primitive.ObjectID      `json:"id" bson:"_id,omitempty"`
	UUID      string                  `json:"uuid" bson:"uuid"`
	Name      string                  `json:"name" bson:"name"`
	Body      string                  `json:"body" bson:"body"`
	Status    enum.NotificationStatus `json:"type" bson:"type"`
	IsRead    bool                    `json:"is_read" bson:"is_read"`
	CreatedAt time.Time               `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time               `json:"updated_at" bson:"updated_at"`
}
