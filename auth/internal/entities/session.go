package entities

import "time"

type NewSession struct {
	AuthID    int64
	Token     string
	UserAgent string
	IPAddress string
	ExpiresAt time.Time
}
