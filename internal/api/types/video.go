package types

import (
	"go-fitness/internal/api/enum"
	"time"
)

type Video struct {
	ID          int64            `json:"id"`
	UUID        string           `json:"uuid"`
	Name        string           `json:"name"`
	HashName    string           `json:"hash_name"`
	Description string           `json:"description,omitempty"`
	Status      enum.VideoStatus `json:"status"`
	Duration    float64          `json:"duration"`
	DeletedAt   *time.Time       `json:"deleted_at"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type VideoPosition struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	VideoID   int64     `json:"video_id"`
	Position  float64   `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
