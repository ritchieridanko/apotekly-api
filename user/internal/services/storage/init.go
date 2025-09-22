package storage

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/ritchieridanko/apotekly-api/user/internal/entities"
	"go.opentelemetry.io/otel"
)

const storageErrorTracer string = "service.storage"

type StorageService interface {
	Upload(ctx context.Context, data *entities.NewUpload) (result *uploader.UploadResult, err error)
}

type storageService struct {
	instance *cloudinary.Cloudinary
}

func NewService(instance *cloudinary.Cloudinary) StorageService {
	return &storageService{instance}
}

func (ss *storageService) Upload(ctx context.Context, data *entities.NewUpload) (*uploader.UploadResult, error) {
	ctx, span := otel.Tracer(storageErrorTracer).Start(ctx, "Upload")
	defer span.End()

	params := uploader.UploadParams{
		PublicID:       data.PublicID,
		PublicIDPrefix: data.PublicIDPrefix,
		Folder:         data.Folder,
		Overwrite:      data.Overwrite,
	}

	return ss.instance.Upload.Upload(ctx, data.File, params)
}
