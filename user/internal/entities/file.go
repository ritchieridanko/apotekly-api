package entities

import "io"

type NewUpload struct {
	File           io.Reader
	PublicID       string
	PublicIDPrefix string
	Folder         string
	Overwrite      *bool
}
