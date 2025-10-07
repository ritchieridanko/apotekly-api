package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
)

func NewInstance(bucket, key, secret string) (*cloudinary.Cloudinary, error) {
	storage, err := cloudinary.NewFromParams(bucket, key, secret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	_, err = storage.Admin.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping storage: %w", err)
	}

	log.Println("✅ initialized storage")
	return storage, nil
}
