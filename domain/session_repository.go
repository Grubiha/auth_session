package domain

import (
	"context"
	"time"
)

type SessionRepository interface {
	// SD + RAM
	Create(ctx context.Context, dto CreateSessionDto, ttl, refreshTtl time.Duration) (string, error)
	GetUserSessionCount(ctx context.Context, dto FindSessionWithRoleDto) (int, error)
	DeleteOldestUserSession(ctx context.Context, dto FindSessionWithRoleDto) error
	Delete(ctx context.Context, dto FindSessionDto) error

	// RAM only
	FindSessionInfo(ctx context.Context, dto FindSessionDto) (SessionInfo, error)
}
