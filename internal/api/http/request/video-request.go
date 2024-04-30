package request

import "mime/multipart"

type VideoUploadRequest struct {
	File        multipart.FileHeader `json:"file" validate:"required"`
	Name        string               `json:"name" validate:"required"`
	Description string               `json:"description"`
}

type VideoSavePositionRequest struct {
	Position float64 `json:"position" validate:"required"`
}

type VideoUpdateRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}
