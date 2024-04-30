package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"go-fitness/external/config"
	"go-fitness/external/logger/sl"
	"go-fitness/internal/api/data"
	"go-fitness/internal/api/enum"
	"go-fitness/internal/api/repository"
	"go-fitness/internal/api/types"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type VideoService struct {
	log                 *slog.Logger
	cfg                 *config.Config
	notificationService NotificationServiceInterface
	videoRepo           repository.VideoRepositoryInterface
	transcodeQueue      VideoTranscodeTaskChan
}

type UploadAndTranscodeQueueInterface interface {
	WaitForTranscodeVideoSignals()
}

type VideoServiceInterface interface {
	ProcessUpload(context.Context, data.VideoUploadData) error
	ProcessGetVideoPlayListByUUID(context.Context, string) ([]byte, error)
	ProcessGetVideoTS(context.Context, string) ([]byte, error)
	ProcessGetVideoM3U8(string) ([]byte, error)
	ProcessDeleteVideo(context.Context, string) error
	ProcessUpdateVideoInfo(context.Context, types.Video) error
	ProcessGetVideoPosition(context.Context, int64, string) (float64, error)
	ProcessSaveOrUpdateVideoPosition(context.Context, int64, string, float64) error
	ProcessSoftDeleteVideo(context.Context, string) error

	ProcessGetVideoList(context.Context) ([]VideoResponse, error)
	ProcessGetVideoListWithPosition(context.Context, int64, map[string]interface{}) ([]VideoResponse, error)
}

func NewVideoService(
	log *slog.Logger,
	cfg *config.Config,
	notificationService NotificationServiceInterface,
	videoRepo repository.VideoRepositoryInterface,
	transcodeQueue VideoTranscodeTaskChan,
) *VideoService {
	return &VideoService{
		log:                 log,
		cfg:                 cfg,
		notificationService: notificationService,
		videoRepo:           videoRepo,
		transcodeQueue:      transcodeQueue,
	}
}

// VideoResponse is a response body for Video data
// @Description Video data along with status
// @Accept json
// @Produce json
type VideoResponse struct {
	UUID        string  `json:"uuid"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Status      string  `json:"status,omitempty"`
	Duration    float64 `json:"duration"`

	Position *float64 `json:"position,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type VideoTranscodeTask struct {
	UploadPath string
	VideoID    int64
	DstPath    string
	ChunkHash  string
}

type VideoTranscodeTaskChan chan VideoTranscodeTask

func NewVideoTranscodeTask() VideoTranscodeTaskChan {
	return make(VideoTranscodeTaskChan, 50)
}

func (s *VideoService) WaitForTranscodeVideoSignals() {
	const op = "VideoService.WaitForTranscodeVideoSignals"

	// Quanta of work to be done by each worker
	workerCount := s.cfg.VideoService.TranscodeVideoWorkerCount

	log := s.log.With(sl.String("op", op))

	log.Info("Initializing worker pool")

	for i := 0; i < workerCount; i++ {
		go s.transcodeVideoWorker(i)
	}
}

// transcodeVideoWorker is a method to transcode video worker
func (s *VideoService) transcodeVideoWorker(workerID int) {
	const op = "VideoService.uploadVideoWorker"

	log := s.log.With(
		sl.String("op", op),
		sl.Int("workerID", workerID),
	)

	log.Info("Worker started")

	for task := range s.transcodeQueue {
		log.Info("Processing upload", sl.Any("task", task))
		if err := s.processTranscode(context.Background(), task.UploadPath, task.VideoID, task.DstPath, task.ChunkHash); err != nil {
			log.Error("Failed to process upload", sl.Err(err))
			return
		}
	}

	log.Info("Worker stopped")
}

// ProcessSaveOrUpdateVideoPosition is a method to process save or update video position
func (s *VideoService) ProcessSaveOrUpdateVideoPosition(ctx context.Context, userID int64, videoUUID string, position float64) error {
	const op = "VideoService.ProcessSaveVideoPosition"

	log := s.log.With(
		sl.String("op", op),
	)

	video, err := s.videoRepo.GetByUUID(ctx, videoUUID)
	if err != nil {
		log.Error("failed to get video by uuid", sl.Err(err))
		return errors.New("failed to get video by uuid")
	}

	videoPosition := types.VideoPosition{
		UserID:   userID,
		VideoID:  video.ID,
		Position: position,
	}

	p, err := s.videoRepo.GetVideoPositionByIDAndUserID(ctx, video.ID, userID)
	if err != nil {
		if err := s.videoRepo.SaveVideoPosition(ctx, videoPosition); err != nil {
			log.Error("failed to save video position", sl.Err(err))
			return errors.New("failed to save video position")
		}
		return nil
	}

	p.Position = videoPosition.Position
	if err := s.videoRepo.UpdateVideoPosition(ctx, *p); err != nil {
		log.Error("failed to update video position", sl.Err(err))
		return errors.New("failed to update video position")
	}

	return nil
}

// readFile is a method to read file from the storage path
func (s *VideoService) readFile(path string) ([]byte, error) {
	const op = "VideoService.readFile"

	log := s.log.With(
		sl.String("op", op),
		sl.String("path", path),
	)

	slurp, err := os.ReadFile(path)
	if err != nil {
		log.Error("failed to read file", sl.Err(err))
		return nil, errors.New("failed to read file")
	}
	return slurp, nil
}

// ProcessGetVideoPlayListByUUID is a method to process video playlist by UUID and return the video file
func (s *VideoService) ProcessGetVideoPlayListByUUID(ctx context.Context, uuid string) ([]byte, error) {
	const op = "VideoService.ProcessGetVideoM3U8ByUUID"

	log := s.log.With(
		sl.String("op", op),
		sl.String("uuid", uuid),
	)

	video, err := s.videoRepo.GetByUUID(ctx, uuid)
	if err != nil {
		log.Error("failed to get video by uuid", sl.Err(err))
		return nil, errors.New("failed to get video by uuid")
	}

	videoPath := fmt.Sprintf("%s/%s/%s/playlist.m3u8", s.cfg.HTTPServer.StoragePath, s.cfg.VideoService.VideoPath, video.HashName)
	return s.readFile(videoPath)
}

// ProcessGetVideoM3U8 is a method to process video M3U8 and return the video file
func (s *VideoService) ProcessGetVideoM3U8(url string) ([]byte, error) {
	const op = "VideoService.ProcessGetVideoM3U8"

	log := s.log.With(
		sl.String("op", op),
		sl.String("url", url),
	)

	var lastError error

	hashName, resolution, err := s.parseURL(url)
	if err != nil {
		log.Error("failed to parse hash", sl.Err(err))
		return nil, errors.New("failed to parse hash")
	}

	videoPath := fmt.Sprintf("%s/%s/%s/%s", s.cfg.HTTPServer.StoragePath, s.cfg.VideoService.VideoPath, hashName, resolution)

	video, err := s.readFile(videoPath)
	if err != nil {
		log.Error("failed to read file", sl.Err(err))

		var resolutions = s.cfg.VideoService.Resolutions

		for _, res := range resolutions {
			if res == resolution {
				continue
			}

			res = fmt.Sprintf("%s.m3u8", res)
			videoPath = fmt.Sprintf("%s/%s/%s/%s", s.cfg.HTTPServer.StoragePath, s.cfg.VideoService.VideoPath, hashName, res)
			video, err = s.readFile(videoPath)
			if err == nil {
				return video, nil
			}
			lastError = err
		}

		log.Error("failed to find any suitable video file", sl.Err(lastError))
		return nil, lastError
	}

	return video, nil
}

// ProcessGetVideoTS is a method to process video TS and return the video file
func (s *VideoService) ProcessGetVideoTS(ctx context.Context, url string) ([]byte, error) {
	const op = "VideoService.ProcessGetVideoTS"

	log := s.log.With(
		sl.String("op", op),
		sl.String("url", url),
	)

	hashName, resolution, err := s.parseURL(url)
	if err != nil {
		log.Error("failed to parse hash", sl.Err(err))
		return nil, errors.New("failed to parse hash")
	}

	videoPath := fmt.Sprintf("%s/%s/%s/%s", s.cfg.HTTPServer.StoragePath, s.cfg.VideoService.VideoPath, hashName, resolution)
	return s.readFile(videoPath)
}

// parseURL is a method to parse the URL and return the hash and resolution
func (s *VideoService) parseURL(pathstr string) (string, string, error) {
	paths := strings.SplitN(strings.TrimLeft(pathstr, "/"), "/", -1)
	if len(paths) < 2 {
		return "", "", fmt.Errorf("invalid path: %s", pathstr)
	}
	return paths[len(paths)-2], paths[len(paths)-1], nil
}

// ProcessUpload is a method to process video upload and store it in the storage path
func (s *VideoService) ProcessUpload(
	ctx context.Context,
	data data.VideoUploadData,
) error {
	const op = "VideoService.ProcessUpload"

	log := s.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		log.Error("context done")
		return errors.New("context done")
	default:
		log.Info("processing upload")
	}

	dstPath, uploadPath, chunkHash, err := s.uploadFile(data)
	if err != nil {
		log.Error("failed to upload file", sl.Err(err))
		return errors.New("failed to upload file")
	}

	duration, err := s.getVideoDuration(dstPath)
	if err != nil {
		log.Error("failed to get video duration", sl.Err(err))
		return errors.New("failed to get video duration")
	}

	video := types.Video{
		Name:        data.Name,
		Description: data.Description,
		HashName:    chunkHash,
		Status:      enum.VideoStatusProcessing,
		Duration:    duration,
	}

	videoID, err := s.videoRepo.Create(ctx, video)
	if err != nil {
		log.Error("failed to create video", sl.Err(err))
		return errors.New("failed to create video")
	}

	videoTranscodeTask := VideoTranscodeTask{
		UploadPath: uploadPath,
		VideoID:    videoID,
		DstPath:    dstPath,
		ChunkHash:  chunkHash,
	}

	s.transcodeQueue <- videoTranscodeTask

	return nil
}

func (s *VideoService) getVideoDuration(filePath string) (float64, error) {
	const op = "VideoService.getVideoDuration"

	log := s.log.With(
		sl.String("op", op),
		sl.String("file_path", filePath),
	)

	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", filePath)
	out, err := cmd.Output()
	if err != nil {
		log.Error("error running ffprobe", sl.Err(err))
		return 0, errors.New("error running ffprobe")
	}
	durationStr := strings.TrimSpace(string(out))
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		log.Error("error converting duration", sl.Err(err))
		return 0, errors.New("error converting duration")
	}
	return duration, nil
}

// uploadFile is a method to upload file to the storage path
func (s *VideoService) uploadFile(data data.VideoUploadData) (string, string, string, error) {
	const op = "VideoService.uploadFile"

	log := s.log.With(
		sl.String("op", op),
	)

	chunkHash := _hash(fmt.Sprintf("%s-%d", data.Header.Filename, data.Header.Size))

	uploadPath := fmt.Sprintf("%s/%s/%s", s.cfg.HTTPServer.StoragePath, s.cfg.VideoService.VideoPath, chunkHash)

	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		log.Error("failed to create upload directory", sl.Err(err))
		return "", "", "", errors.New("failed to create upload directory")
	}

	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		err := os.Mkdir(uploadPath, os.ModePerm)
		if err != nil {
			log.Error("failed to create upload directory", sl.Err(err))
			return "", "", "", errors.New("failed to create upload directory")
		}
	}

	dstPath := filepath.Join(uploadPath, data.Header.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		log.Error("failed to create file", sl.Err(err))
		return "", "", "", errors.New("failed to create file")
	}
	defer func(dst *os.File) {
		if err := dst.Close(); err != nil {
			log.Error("failed to close file", sl.Err(err))
		}
	}(dst)

	if _, err = io.Copy(dst, data.File); err != nil {
		log.Error("failed to copy file", sl.Err(err))
		return "", "", "", errors.New("failed to copy file")
	}

	return dstPath, uploadPath, chunkHash, nil
}

// processTranscode is a method to process video transcoding and chunking
func (s *VideoService) processTranscode(
	ctx context.Context,
	uploadPath string,
	videoID int64,
	dstPath string,
	chunkHash string,
) error {
	const op = "VideoService.processTranscode"

	log := s.log.With(
		sl.String("op", op),
	)

	uploadSuccessful := false

	defer func() {
		//TODO translate to romanian
		notification := types.Notification{
			Name:   "Upload Status",
			Body:   "The upload was successful.",
			Status: enum.NotificationStatusSuccess,
			IsRead: false,
		}

		if !uploadSuccessful {
			notification.Body = "The upload has failed."
			notification.Status = enum.NotificationStatusError

			if updateErr := s.videoRepo.UpdateStatus(ctx, videoID, enum.VideoStatusFailed); updateErr != nil {
				log.Error("failed to update video status to failed", sl.Err(updateErr))
			}

			if rmErr := os.RemoveAll(uploadPath); rmErr != nil {
				log.Error("failed to remove folder", sl.Err(rmErr))
			}
		} else {
			if rmErr := os.Remove(dstPath); rmErr != nil {
				log.Error("failed to remove file after failed upload", sl.Err(rmErr))
			}
		}

		if err := s.notificationService.ProcessNotification(ctx, notification); err != nil {
			log.Error("failed to send notification", sl.Err(err))
		}
	}()

	if err := s.transcodeAndChunk(uploadPath, dstPath); err != nil {
		log.Error("failed to transcode and chunk video", sl.Err(err))
		return errors.New("failed to transcode and chunk video")
	}

	if err := s.createMasterM8U3PlayList(uploadPath, chunkHash); err != nil {
		log.Error("failed to create master m8u3 playlist", sl.Err(err))
		return errors.New("failed to create master m8u3 playlist")
	}

	if err := s.videoRepo.UpdateStatus(ctx, videoID, enum.VideoStatusProcessed); err != nil {
		log.Error("failed to update video status to processed", sl.Err(err))
		return errors.New("failed to update video status to processed")
	}

	uploadSuccessful = true

	return nil
}

// transcodeAndChunk is a method to transcode and chunk video into smaller segments using ffmpeg
func (s *VideoService) transcodeAndChunk(uploadPath, videoPath string) error {
	const op = "VideoService.transcodeAndChunk"

	log := s.log.With(
		sl.String("op", op),
		sl.String("upload_path", uploadPath),
		sl.String("video_path", videoPath),
	)

	log.Info("transcoding and chunking video")

	var resolutions []string
	for res := range s.cfg.VideoService.Resolutions {
		resolutions = append(resolutions, res)
	}

	sort.Slice(resolutions, func(i, j int) bool {
		resI := strings.Split(resolutions[i], "x")
		resJ := strings.Split(resolutions[j], "x")
		if len(resI) < 2 || len(resJ) < 2 {
			return false
		}
		heightI, _ := strconv.Atoi(resI[1])
		heightJ, _ := strconv.Atoi(resJ[1])
		return heightI < heightJ
	})

	for _, res := range resolutions {
		label := s.cfg.VideoService.Resolutions[res]
		outputPath := fmt.Sprintf("%s/%s.m3u8", uploadPath, label)
		segmentFilename := fmt.Sprintf("%s/%s_%%03d.ts", uploadPath, label)
		cmd := exec.Command("ffmpeg", "-i", videoPath,
			"-profile:v", "baseline", "-level", "3.0",
			"-s", res,
			"-start_number", "0",
			"-hls_time", "10",
			"-hls_list_size", "0",
			"-f", "hls",
			"-hls_segment_filename",
			segmentFilename,
			outputPath)
		if err := cmd.Run(); err != nil {
			log.Error("failed to transcode video", sl.String("resolution", res), sl.Err(err))
			return errors.New("failed to transcode video for resolution " + label)
		}
		log.Info("transcoded video", sl.String("resolution", label))
	}

	return nil
}

// createMasterM8U3PlayList is a method to create master m8u3 playlist
func (s *VideoService) createMasterM8U3PlayList(uploadPath string, chunkHash string) error {
	const op = "VideoService.createMasterM8U3PlayList"

	log := s.log.With(
		sl.String("op", op),
		sl.String("upload_path", uploadPath),
	)

	log.Info("creating master m8u3 playlist")

	masterM8U3PlayListPath := fmt.Sprintf("%s/%s.m3u8", uploadPath, "playlist")

	masterM8U3PlayList, err := os.Create(masterM8U3PlayListPath)
	if err != nil {
		log.Error("failed to create master m8u3 playlist", sl.Err(err))
		return errors.New("failed to create master m8u3 playlist")
	}
	defer func(masterM8U3PlayList *os.File) {
		if err := masterM8U3PlayList.Close(); err != nil {
			log.Error("failed to close master m8u3 playlist", sl.Err(err))
		}
	}(masterM8U3PlayList)

	var buffer bytes.Buffer
	buffer.WriteString("#EXTM3U\n")
	buffer.WriteString("#EXT-X-VERSION:3\n")
	buffer.WriteString("#EXT-X-STREAM-INF:BANDWIDTH=800000,RESOLUTION=640x360\n")
	buffer.WriteString(chunkHash + "/360.m3u8\n")
	buffer.WriteString("#EXT-X-STREAM-INF:BANDWIDTH=1400000,RESOLUTION=854x480\n")
	buffer.WriteString(chunkHash + "/480.m3u8\n")
	buffer.WriteString("#EXT-X-STREAM-INF:BANDWIDTH=2800000,RESOLUTION=1280x720\n")
	buffer.WriteString(chunkHash + "/720.m3u8\n")
	buffer.WriteString("#EXT-X-STREAM-INF:BANDWIDTH=5000000,RESOLUTION=1920x1080\n")
	buffer.WriteString(chunkHash + "/1080.m3u8\n")

	if _, err := masterM8U3PlayList.Write(buffer.Bytes()); err != nil {
		log.Error("failed to write master m8u3 playlist", sl.Err(err))
		return errors.New("failed to write master m8u3 playlist")
	}

	return nil
}

// getBitrate is a method to get the bitrate of the video
func (s *VideoService) getBitrate(filePath string) (int64, error) {
	const op = "VideoService.getBitrate"

	log := s.log.With(
		sl.String("op", op),
	)
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=bit_rate", "-of", "default=noprint_wrappers=1:nokey=1", filePath)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Error("error running ffprobe", sl.Err(err), sl.String("stderr", stderr.String()))
		return 0, errors.New("error running ffprobe")
	}

	bitrate, err := strconv.Atoi(strings.TrimSpace(stdout.String()))
	if err != nil {
		log.Error("error converting bitrate", sl.Err(err))
		return 0, errors.New("error converting bitrate")
	}

	return int64(bitrate), err
}

// ProcessDeleteVideo is a method to process video deletion
func (s *VideoService) ProcessDeleteVideo(ctx context.Context, uuid string) error {
	const op = "VideoService.ProcessDeleteVideo"

	log := s.log.With(
		sl.String("op", op),
		sl.String("uuid", uuid),
	)

	video, err := s.videoRepo.GetByUUID(ctx, uuid)
	if err != nil {
		log.Error("failed to get video by uuid", sl.Err(err))
		return errors.New("failed to get video by uuid")
	}

	videoPath := fmt.Sprintf("%s/%s/%s", s.cfg.HTTPServer.StoragePath, s.cfg.VideoService.VideoPath, video.HashName)
	if err := os.RemoveAll(videoPath); err != nil {
		log.Error("failed to remove video folder", sl.Err(err))
		return errors.New("failed to remove video folder")
	}

	if err := s.videoRepo.Delete(ctx, video.ID); err != nil {
		log.Error("failed to delete video", sl.Err(err))
		return errors.New("failed to delete video")
	}

	return nil
}

// ProcessSoftDeleteVideo is a method to process soft delete video
func (s *VideoService) ProcessSoftDeleteVideo(ctx context.Context, uuid string) error {
	const op = "VideoService.ProcessSoftDeleteVideo"

	log := s.log.With(
		sl.String("op", op),
		sl.String("uuid", uuid),
	)

	video, err := s.videoRepo.GetByUUID(ctx, uuid)
	if err != nil {
		log.Error("failed to get video by uuid", sl.Err(err))
		return errors.New("failed to get video by uuid")
	}

	if err := s.videoRepo.SoftDelete(ctx, video.ID); err != nil {
		log.Error("failed to soft delete video", sl.Err(err))
		return errors.New("failed to soft delete video")
	}

	return nil
}

// ProcessGetVideoList is a method to process getting video list
func (s *VideoService) ProcessGetVideoList(ctx context.Context) ([]VideoResponse, error) {
	const op = "VideoService.ProcessGetVideoList"

	log := s.log.With(
		sl.String("op", op),
	)

	videos, err := s.videoRepo.GetList(ctx, nil)
	if err != nil {
		log.Error("failed to get videos", sl.Err(err))
		return nil, errors.New("failed to get videos")
	}

	var response []VideoResponse

	for _, video := range videos {
		resp := VideoResponse{
			UUID:        video.UUID,
			Name:        video.Name,
			Description: video.Description,
			Status:      video.Status.String(),
			Duration:    video.Duration,
			CreatedAt:   video.CreatedAt,
			UpdatedAt:   video.UpdatedAt,
		}

		response = append(response, resp)
	}

	return response, nil
}

// ProcessGetVideoListWithPosition is a method to process getting video list with position
func (s *VideoService) ProcessGetVideoListWithPosition(
	ctx context.Context,
	userID int64,
	filter map[string]interface{},
) ([]VideoResponse, error) {
	const op = "VideoService.ProcessGetVideoList"

	log := s.log.With(
		sl.String("op", op),
	)

	filter = map[string]interface{}{
		"status": enum.VideoStatusProcessed,
	}

	videos, err := s.videoRepo.GetList(ctx, filter)
	if err != nil {
		log.Error("failed to get videos", sl.Err(err))
		return nil, errors.New("failed to get videos")
	}

	var response []VideoResponse
	for _, video := range videos {
		resp := VideoResponse{
			UUID:        video.UUID,
			Name:        video.Name,
			Description: video.Description,
			Duration:    video.Duration,
			CreatedAt:   video.CreatedAt,
		}

		videoPosition, err := s.videoRepo.GetVideoPositionByIDAndUserID(ctx, video.ID, userID)
		if err == nil {
			resp.Position = &videoPosition.Position
		}

		response = append(response, resp)
	}

	return response, nil
}

// ProcessUpdateVideoInfo is a method to process updating video info
func (s *VideoService) ProcessUpdateVideoInfo(ctx context.Context, video types.Video) error {
	const op = "VideoService.ProcessUpdateVideoInfo"

	log := s.log.With(
		sl.String("op", op),
	)

	if err := s.videoRepo.Update(ctx, video); err != nil {
		log.Error("failed to update video", sl.Err(err))
		return errors.New("failed to update video")
	}

	return nil
}

// ProcessGetVideoPosition is a method to process getting video position
func (s *VideoService) ProcessGetVideoPosition(ctx context.Context, userID int64, videoUUID string) (float64, error) {
	const op = "VideoService.ProcessGetVideoPosition"

	log := s.log.With(
		sl.String("op", op),
	)

	video, err := s.videoRepo.GetByUUID(ctx, videoUUID)
	if err != nil {
		log.Error("failed to get video by uuid", sl.Err(err))
		return 0, errors.New("failed to get video by uuid")
	}

	position, err := s.videoRepo.GetVideoPositionByIDAndUserID(ctx, video.ID, userID)
	if err != nil {
		log.Error("failed to get video by uuid", sl.Err(err))
		return 0, errors.New("failed to get video by uuid")
	}

	return position.Position, nil
}

func _hash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}
