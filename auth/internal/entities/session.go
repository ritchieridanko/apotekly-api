package entities

import "time"

type NewSession struct {
	AuthID    int64
	Token     string
	UserAgent string
	IPAddress string
	ExpiresAt time.Time
}

type ReissueSession struct {
	AuthID    int64
	ParentID  int64
	Token     string
	UserAgent string
	IPAddress string
	ExpiresAt time.Time
}
