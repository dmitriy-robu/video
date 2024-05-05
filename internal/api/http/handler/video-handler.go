package handler

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"go-fitness/external/logger/sl"
	"go-fitness/external/response"
	"go-fitness/external/validation"
	"go-fitness/internal/api/data"
	"go-fitness/internal/api/http/request"
	"go-fitness/internal/api/service"
	"go-fitness/internal/api/types"
	"go.mongodb.org/mongo-driver/bson"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type VideoHandler struct {
	log          *slog.Logger
	videoService service.VideoServiceInterface
	validation   *validator.Validate
}

func NewVideoHandler(
	log *slog.Logger,
	videoService service.VideoServiceInterface,
) *VideoHandler {
	return &VideoHandler{
		log:          log,
		videoService: videoService,
		validation:   validator.New(),
	}
}

func (h *VideoHandler) GetVideos() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op string = "VideoHandler.GetVideos"

		log := h.log.With(
			sl.String("op", op),
		)

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		videos, err := h.videoService.ProcessGetVideoList(ctx)
		if err != nil {
			log.Error("failed to get form file", sl.Err(err))
			response.Respond(w, response.Response{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
				Data:    err.Error(),
			})
			return
		}

		response.Respond(w, response.Response{
			Status:  http.StatusOK,
			Message: "ok",
			Data:    videos,
		})
		return
	}
}

func (h *VideoHandler) GetVideosWithPositions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op string = "VideoHandler.GetVideosWithPositions"

		log := h.log.With(
			sl.String("op", op),
		)

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		userID := ctx.Value("user").(types.User).ID

		filter := bson.M{}

		videos, err := h.videoService.ProcessGetVideoListWithPosition(ctx, userID, filter)
		if err != nil {
			log.Error("failed to get form file", sl.Err(err))
			response.Respond(w, response.Response{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
				Data:    err.Error(),
			})
			return
		}

		response.Respond(w, response.Response{
			Status:  http.StatusOK,
			Message: "ok",
			Data:    videos,
		})

		return
	}
}

func (h *VideoHandler) GetVideo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op string = "VideoHandler.GetVideos"

		log := h.log.With(
			sl.String("op", op),
		)

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if strings.Contains(r.URL.Path, ".ts") {
			video, err := h.videoService.ProcessGetVideoTS(ctx, r.URL.Path)
			if err != nil {
				log.Error("failed to get videos", sl.Err(err))
				response.Respond(w, response.Response{
					Status:  http.StatusInternalServerError,
					Message: "internal server error",
					Data:    err.Error(),
				})
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "video/mp2ts")
			w.Header().Set("Content-Length", strconv.Itoa(len(video)))

			_, err = w.Write(video)
			if err != nil {
				log.Error("failed to write video", sl.Err(err))
				response.Respond(w, response.Response{
					Status:  http.StatusInternalServerError,
					Message: "internal server error",
					Data:    err.Error(),
				})
				return
			}
		} else if strings.Contains(r.URL.Path, ".m3u8") {
			video, err := h.videoService.ProcessGetVideoM3U8(r.URL.Path)
			if err != nil {
				log.Error("failed to get videos", sl.Err(err))
				response.Respond(w, response.Response{
					Status:  http.StatusInternalServerError,
					Message: "internal server error",
					Data:    err.Error(),
				})
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Content-Type", "application/x-mpegURL")

			_, err = w.Write(video)
			if err != nil {
				log.Error("failed to write video", sl.Err(err))
				response.Respond(w, response.Response{
					Status:  http.StatusInternalServerError,
					Message: "internal server error",
					Data:    err.Error(),
				})
				return
			}
		} else {
			videoUUID := chi.URLParam(r, "uuid")

			if videoUUID == "" {
				log.Error("uuid is required")
				response.Respond(w, response.Response{
					Status:  http.StatusInternalServerError,
					Message: "internal server error",
					Data:    "UUID is required",
				})
				return
			}

			video, err := h.videoService.ProcessGetVideoPlayListByUUID(ctx, videoUUID)
			if err != nil {
				log.Error("failed to get videos", sl.Err(err))
				response.Respond(w, response.Response{
					Status:  http.StatusInternalServerError,
					Message: "internal server error",
					Data:    err.Error(),
				})
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Content-Type", "application/x-mpegURL")

			_, writeErr := w.Write(video)
			if writeErr != nil {
				log.Error("failed to write video", sl.Err(writeErr))
				response.Respond(w, response.Response{
					Status:  http.StatusInternalServerError,
					Message: "internal server error",
					Data:    err.Error(),
				})
				return
			}
		}
	}
}

// GetVideoPosition gets the video position by uuid and user uuid
func (h *VideoHandler) GetVideoPosition() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op string = "VideoHandler.GetVideoPosition"

		log := h.log.With(
			sl.String("op", op),
		)

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		videoUUID := chi.URLParam(r, "uuid")
		if videoUUID == "" {
			log.Error("uuid is required")
			response.Respond(w, response.Response{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
				Data:    "uuid is required",
			})
			return
		}

		userID := ctx.Value("user").(types.User).ID

		position, err := h.videoService.ProcessGetVideoPosition(ctx, userID, videoUUID)
		if err != nil {
			log.Error("failed to get video position", sl.Err(err))
			response.Respond(w, response.Response{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
				Data:    err,
			})
			return
		}

		response.Respond(w, response.Response{
			Status:  http.StatusOK,
			Message: "ok",
			Data:    position,
		})

		return
	}
}

// SaveVideoPosition saves the video position to the database
func (h *VideoHandler) SaveVideoPosition() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op string = "VideoHandler.SaveVideoPosition"

		log := h.log.With(
			sl.String("op", op),
		)

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		videoUUID := chi.URLParam(r, "uuid")
		if videoUUID == "" {
			log.Error("video_uuid is required")
			response.Respond(w, response.Response{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
				Data:    "video_uuid is required",
			})
			return
		}

		var positionRequest request.VideoSavePositionRequest

		if err := render.DecodeJSON(r.Body, &positionRequest); err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			response.Respond(w, response.Response{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
				Data:    err.Error(),
			})
			return
		}

		var validateErr validator.ValidationErrors
		if err := h.validation.Struct(positionRequest); err != nil {
			errors.As(err, &validateErr)
			log.Error("invalid request", sl.Err(validateErr))
			response.Respond(w, response.Response{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
				Data:    validation.ValidationError(validateErr).Error(),
			})
			return
		}

		userID := ctx.Value("user").(types.User).ID

		if err := h.videoService.ProcessSaveOrUpdateVideoPosition(ctx, userID, videoUUID, positionRequest.Position); err != nil {
			log.Error("failed to save video position", sl.Err(err))
			response.Respond(w, response.Response{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
				Data:    err.Error(),
			})
			return
		}

		response.Respond(w, response.Response{
			Status:  http.StatusOK,
			Message: "ok",
			Data:    "ok",
		})

		return
	}
}

// ProcessUpload processes the video upload
// It returns a http.HandlerFunc
func (h *VideoHandler) ProcessUpload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op string = "VideoHandler.ProcessUpload"

		log := h.log.With(
			sl.String("op", op),
		)

		log.Info("processing upload")

		ctx := r.Context()

		// 8 << 30 is 8GB
		if err := r.ParseMultipartForm(8 << 30); err != nil {
			log.Error("failed to parse multipart form", sl.Err(err))
			response.Respond(w, response.Response{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
				Data:    err.Error(),
			})
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			log.Error("failed to get form file", sl.Err(err))
			response.Respond(w, response.Response{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
				Data:    err.Error(),
			})
			return
		}
		defer func(file multipart.File) {
			if err := file.Close(); err != nil {
				log.Error("failed to close file", sl.Err(err))
			}
		}(file)

		uploadRequest := request.VideoUploadRequest{
			File:        *header,
			Name:        r.FormValue("name"),
			Description: r.FormValue("description"),
		}

		var validateErr validator.ValidationErrors
		if err := h.validation.Struct(uploadRequest); err != nil {
			errors.As(err, &validateErr)
			log.Error("invalid request", sl.Err(validateErr))
			response.Respond(w, response.Response{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
				Data:    validation.ValidationError(validateErr).Error(),
			})
			return
		}

		uploadData := data.VideoUploadData{
			File:        file,
			Header:      header,
			Name:        uploadRequest.Name,
			Description: uploadRequest.Description,
		}

		if err = h.videoService.ProcessUpload(ctx, uploadData); err != nil {
			log.Error("failed to process upload", sl.Err(err))
			response.Respond(w, response.Response{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
				Data:    err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "ok")
		return
	}
}

// DeleteVideo deletes a video by uuid from the database and storage
func (h *VideoHandler) DeleteVideo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op string = "VideoHandler.DeleteVideo"

		log := h.log.With(
			sl.String("op", op),
		)

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		uuid := chi.URLParam(r, "uuid")
		if uuid == "" {
			log.Error("uuid is required")
			response.Respond(w, response.Response{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
				Data:    "uuid is required",
			})
			return
		}

		err := h.videoService.ProcessDeleteVideo(ctx, uuid)
		if err != nil {
			log.Error("failed to delete video", sl.Err(err))
			response.Respond(w, response.Response{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
				Data:    err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "ok")
		return
	}
}

// SoftDeleteVideo soft deletes, update deleted_at field of a video by uuid from the database
func (h *VideoHandler) SoftDeleteVideo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op string = "VideoHandler.DeleteVideo"

		log := h.log.With(
			sl.String("op", op),
		)

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		uuid := chi.URLParam(r, "uuid")
		if uuid == "" {
			log.Error("uuid is required")
			response.Respond(w, response.Response{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
				Data:    "uuid is required",
			})
			return
		}

		err := h.videoService.ProcessSoftDeleteVideo(ctx, uuid)
		if err != nil {
			log.Error("failed to delete video", sl.Err(err))
			response.Respond(w, response.Response{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
				Data:    err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "ok")
		return
	}
}

func (h *VideoHandler) UpdateVideoInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op string = "VideoHandler.UpdateVideoInfo"

		log := h.log.With(
			sl.String("op", op),
		)

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		uuid := chi.URLParam(r, "uuid")
		if uuid == "" {
			log.Error("uuid is required")
			response.Respond(w, response.Response{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
				Data:    "uuid is required",
			})
			return
		}

		var updateRequest request.VideoUpdateRequest

		if err := render.DecodeJSON(r.Body, &updateRequest); err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			response.Respond(w, response.Response{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
				Data:    err.Error(),
			})
			return
		}

		var validateErr validator.ValidationErrors
		if err := h.validation.Struct(updateRequest); err != nil {
			errors.As(err, &validateErr)
			log.Error("invalid request", sl.Err(validateErr))
			response.Respond(w, response.Response{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
				Data:    validation.ValidationError(validateErr).Error(),
			})
			return
		}

		video := types.Video{
			UUID:        uuid,
			Name:        updateRequest.Name,
			Description: updateRequest.Description,
		}

		if err := h.videoService.ProcessUpdateVideoInfo(ctx, video); err != nil {
			log.Error("failed to update video info", sl.Err(err))
			response.Respond(w, response.Response{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
				Data:    err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, "ok")
		return
	}
}
