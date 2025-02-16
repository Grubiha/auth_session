package domain

import "time"

type SessionInfo struct {
	UserId   string
	UserName string
	UserRole string
}

type Session struct {
	Id string

	SessionInfo

	ExpiresAt        time.Time
	RefreshExpiresAt time.Time
}
