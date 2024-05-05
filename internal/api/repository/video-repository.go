package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go-fitness/external/db"
	"go-fitness/internal/api/enum"
	"go-fitness/internal/api/types"
	"time"
)

type VideoRepository struct {
	db db.SqlInterface
}

type VideoRepositoryInterface interface {
	Create(context.Context, types.Video) (int64, error)
	Update(context.Context, types.Video) error
	UpdateStatus(context.Context, int64, enum.VideoStatus) error
	GetByUUID(context.Context, string) (types.Video, error)
	GetList(context.Context, map[string]interface{}) ([]types.Video, error)
	Delete(context.Context, int64) error
	SoftDelete(context.Context, int64) error
	GetVideoPositionByIDAndUserID(context.Context, int64, int64) (types.VideoPosition, error)
	SaveVideoPosition(context.Context, types.VideoPosition) error
	UpdateVideoPosition(context.Context, types.VideoPosition) error
}

func NewVideoRepository(
	db db.SqlInterface,
) *VideoRepository {
	return &VideoRepository{
		db: db,
	}
}

func (r *VideoRepository) SaveVideoPosition(ctx context.Context, videoPosition types.VideoPosition) error {
	const op string = "VideoRepository.SaveVideoPosition"

	const query string = `
		INSERT INTO video_positions 
		    (user_id,video_id,position,created_at,updated_at) 
		VALUES (?,?,?,?,?)
	`

	now := time.Now()

	_, err := r.db.GetExecer().ExecContext(ctx, query, videoPosition.UserID, videoPosition.VideoID, videoPosition.Position, now, now)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *VideoRepository) UpdateVideoPosition(ctx context.Context, videoPosition types.VideoPosition) error {
	const op string = "VideoRepository.UpdateVideoPosition"

	const query string = `
		UPDATE video_positions 
		SET position = ?, updated_at = ? 
		WHERE id = ?
	`

	_, err := r.db.GetExecer().ExecContext(ctx, query, videoPosition.Position, time.Now(), videoPosition.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *VideoRepository) Create(ctx context.Context, video types.Video) (int64, error) {
	const op string = "VideoRepository.Create"

	now := time.Now()

	video.UUID = uuid.New().String()
	video.CreatedAt = now
	video.UpdatedAt = now

	const query string = `
		INSERT INTO videos 
		    (uuid,name,hash_name,description,status,duration,created_at,updated_at) 
		VALUES (?,?,?,?,?,?,?,?)
	`

	inId, err := r.db.GetExecer().ExecContext(ctx, query,
		video.UUID,
		video.Name,
		video.HashName,
		video.Description,
		video.Status,
		video.Duration,
		video.CreatedAt,
		video.UpdatedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := inId.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *VideoRepository) Update(ctx context.Context, video types.Video) error {
	const op string = "VideoRepository.Update"

	const query string = `
		UPDATE videos 
		SET name = $1, hash_name = $2, description = $3, status = $4, updated_at = $5 
		WHERE id = $6
	`

	_, err := r.db.GetExecer().ExecContext(ctx, query,
		video.Name,
		video.HashName,
		video.Description,
		video.Status,
		time.Now(),
		video.ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *VideoRepository) UpdateStatus(ctx context.Context, id int64, status enum.VideoStatus) error {
	const op string = "VideoRepository.UpdateStatus"

	const query string = `
		UPDATE videos 
		SET status = $1 
		WHERE id = $2
	`

	_, err := r.db.GetExecer().ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *VideoRepository) GetByUUID(ctx context.Context, uuid string) (types.Video, error) {
	const op string = "VideoRepository.GetByUUID"

	const query string = `
		SELECT id,uuid,name,hash_name,description,status,created_at,updated_at 
		FROM videos 
		WHERE uuid = ? 
		  AND status = ?
	`

	var video types.Video

	if err := r.db.GetExecer().QueryRowContext(ctx, query, uuid, enum.VideoStatusProcessed).Scan(
		&video.ID,
		&video.UUID,
		&video.Name,
		&video.HashName,
		&video.Description,
		&video.Status,
		&video.CreatedAt,
		&video.UpdatedAt,
	); err != nil {
		return video, fmt.Errorf("%s: %w", op, err)
	}

	return video, nil
}

func (r *VideoRepository) GetList(ctx context.Context, filters map[string]interface{}) ([]types.Video, error) {
	const op string = "VideoRepository.GetList"

	var query = `
		SELECT id,uuid,name,hash_name,description,status,duration,created_at,updated_at 
		FROM videos 
		WHERE deleted_at IS NULL
	`

	if len(filters) > 0 {
		for field, value := range filters {
			query += fmt.Sprintf(" AND %s = '%d'", field, value)
		}
	}

	rows, err := r.db.GetExecer().QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var videos []types.Video
	for rows.Next() {
		var video types.Video

		if err = rows.Scan(
			&video.ID,
			&video.UUID,
			&video.Name,
			&video.HashName,
			&video.Description,
			&video.Status,
			&video.Duration,
			&video.CreatedAt,
			&video.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		videos = append(videos, video)

	}
	return videos, nil
}

func (r *VideoRepository) SoftDelete(ctx context.Context, id int64) error {
	const op string = "VideoRepository.SoftDelete"

	const query string = `
		UPDATE videos 
		SET deleted_at = ? 
		WHERE id = ?
	`

	_, err := r.db.GetExecer().ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *VideoRepository) Delete(ctx context.Context, id int64) error {
	const op string = "VideoRepository.Delete"

	const query string = "DELETE FROM videos WHERE id = ?"

	_, err := r.db.GetExecer().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *VideoRepository) GetVideoPositionByIDAndUserID(
	ctx context.Context,
	videoID,
	userID int64,
) (types.VideoPosition, error) {
	const op string = "VideoRepository.GetVideoPositionByUUIDAndUserUUID"

	const query string = `
		SELECT id,user_id,video_id,position,created_at,updated_at 
		FROM video_positions 
		WHERE video_id = ? 
		  AND user_id = ?
	`

	var videoPosition types.VideoPosition

	if err := r.db.GetExecer().QueryRowContext(ctx, query, videoID, userID).Scan(
		&videoPosition.ID,
		&videoPosition.UserID,
		&videoPosition.VideoID,
		&videoPosition.Position,
		&videoPosition.CreatedAt,
		&videoPosition.UpdatedAt,
	); err != nil {
		return videoPosition, fmt.Errorf("%s: %w", op, err)
	}

	return videoPosition, nil
}
