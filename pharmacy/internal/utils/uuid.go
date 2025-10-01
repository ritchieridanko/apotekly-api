package utils

import google "github.com/google/uuid"

func GenerateRandomUUID() (uuid google.UUID) {
	return google.New()
}
