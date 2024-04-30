package data

import "mime/multipart"

type VideoUploadData struct {
	File        multipart.File
	Header      *multipart.FileHeader
	Name        string
	Description string
}
