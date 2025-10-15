package entities

import "time"

type Session struct {
	ID        int64
	AuthID    int64
	ParentID  *int64
	Token     string
	UserAgent string
	IPAddress string
	CreatedAt time.Time
	ExpiresAt time.Time
	RevokedAt *time.Time
}

type CreateSession struct {
	ParentID  *int64
	Token     string
	UserAgent string
	IPAddress string
	ExpiresAt time.Time
}
