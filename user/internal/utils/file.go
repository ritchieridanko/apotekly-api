package utils

import (
	"io"
	"mime/multipart"
)

func FileProcessImage(file multipart.File) (imageBuf []byte, err error) {
	buf, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if err := ValidateImageFile(buf); err != nil {
		return nil, err
	}

	return buf, nil
}
