package storage

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/ritchieridanko/apotekly-api/user/internal/entities"
)

type Storage struct {
	storage *cloudinary.Cloudinary
}

func NewStorage(storage *cloudinary.Cloudinary) *Storage {
	return &Storage{storage}
}

func (s *Storage) Upload(ctx context.Context, data *entities.UploadParams) (*uploader.UploadResult, error) {
	params := uploader.UploadParams{
		PublicID:       data.PublicID,
		PublicIDPrefix: data.Prefix,
		Folder:         data.Folder,
		Overwrite:      data.Overwrite,
	}

	return s.storage.Upload.Upload(ctx, data.File, params)
}
