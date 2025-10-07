package entities

import "mime/multipart"

type UploadParams struct {
	File      multipart.File
	PublicID  string
	Prefix    string
	Folder    string
	Overwrite *bool
}
