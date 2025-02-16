package repos

import "errors"

var (
	ErrPostgresQueryFailed = errors.New("postgres query failed")
	ErrRedisQueryFailed    = errors.New("redis query failed")
	ErrRoleMistmatch       = errors.New("role mismatch")

	ErrUniqueViolation = errors.New("unique violation")

	ErrUserNotFound    = errors.New("user not found")
	ErrSessionNotFound = errors.New("session not found")
)
