package repos

import (
	"context"
	"errors"
	"time"

	"github.com/Grubiha/auth_session/domain"
	"github.com/jackc/pgx/v5"

	// "github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type SessionRepository struct {
	pgPool      *pgxpool.Pool
	redisClient *redis.Client
}

func NewSessionRepository(pgPool *pgxpool.Pool, redisClient *redis.Client) domain.SessionRepository {
	return &SessionRepository{
		pgPool:      pgPool,
		redisClient: redisClient,
	}
}

func (r *SessionRepository) Create(ctx context.Context, dto domain.CreateSessionDto, ttl time.Duration, refreshTtl time.Duration) (string, error) {
	if err := dto.Validate(); err != nil {
		return "", err
	}

	// Находим пользователя
	query := `SELECT "user_name", "user_role" FROM users WHERE "user_id" = $1`
	var userName, userRole string
	err := r.pgPool.QueryRow(ctx, query, dto.UserId).Scan(&userName, &userRole)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrUserNotFound
		}
		return "", errors.Join(ErrPostgresQueryFailed, err)
	}

	// Проверяем можно ли выдать запрошенную роль
	levels := domain.UserRolesLevel
	if levels[userRole] < levels[dto.SessionRole] {
		return "", ErrRoleMistmatch
	}

	// Открываем транзакцию PostgreSQL
	tx, err := r.pgPool.Begin(ctx)
	if err != nil {
		return "", errors.Join(ErrPostgresQueryFailed, err)
	}
	defer tx.Rollback(ctx)

	// Создаем сессию в PostgreSQL
	expiresAt := time.Now().Add(ttl)
	refreshExpiresAt := time.Now().Add(refreshTtl)
	query = `INSERT INTO sessions ("user_id", "session_role", "expires_at", "refresh_expires_at") VALUES ($1, $2, $3, $4) RETURNING "session_id"`
	var newSessionId string
	err = tx.QueryRow(ctx, query,
		dto.UserId,
		dto.SessionRole,
		expiresAt,
		refreshExpiresAt,
	).Scan(&newSessionId)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23503" {
			return "", ErrUserNotFound
		}
		return "", errors.Join(ErrPostgresQueryFailed, err)
	}

	// Создаем информацию о сессии в Redis
	key := "sessions:" + newSessionId
	_, err = r.redisClient.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, key, map[string]interface{}{
			"user_id":   dto.UserId,
			"user_name": userName,
			"user_role": dto.SessionRole,
		})
		pipe.Expire(ctx, key, ttl)
		return nil
	})
	if err != nil {
		return "", errors.Join(ErrRedisQueryFailed, err)
	}

	// Подтверждаем транзакцию PostgreSQL
	if err := tx.Commit(ctx); err != nil {
		return "", errors.Join(ErrPostgresQueryFailed, err)
	}

	return newSessionId, nil
}

func (r *SessionRepository) GetUserSessionCount(ctx context.Context, dto domain.FindSessionWithRoleDto) (int, error) {
	if err := dto.Validate(); err != nil {
		return 0, err
	}
	query := `SELECT COUNT(*) FROM sessions WHERE "user_id" = $1 AND "session_role" = $2 AND "refresh_expires_at" > $3`
	var count int
	err := r.pgPool.QueryRow(ctx, query, dto.Id, dto.SessionRole, time.Now()).Scan(&count)
	if err != nil {
		return 0, errors.Join(ErrPostgresQueryFailed, err)
	}
	return count, nil
}

func (r *SessionRepository) DeleteOldestUserSession(ctx context.Context, dto domain.FindSessionWithRoleDto) error {
	query := `SELECT "session_id" FROM sessions WHERE "user_id" = $1 AND "session_role" = $2 ORDER BY "refresh_expires_at" ASC LIMIT 1`
	var sessionId string
	err := r.pgPool.QueryRow(ctx, query, dto.Id, dto.SessionRole).Scan(&sessionId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return errors.Join(ErrPostgresQueryFailed, err)
	}
	r.Delete(ctx, domain.FindSessionDto{Id: sessionId})
	return nil
}

func (r *SessionRepository) Delete(ctx context.Context, dto domain.FindSessionDto) error {
	if err := dto.Validate(); err != nil {
		return err
	}
	key := "sessions:" + dto.Id
	err := r.redisClient.Del(ctx, key).Err()
	if err != nil {
		return errors.Join(ErrRedisQueryFailed, err)
	}
	query := `DELETE FROM sessions WHERE "session_id" = $1`
	_, err = r.pgPool.Exec(ctx, query, dto.Id)
	if err != nil {
		return errors.Join(ErrPostgresQueryFailed, err)
	}
	return nil
}

func (r *SessionRepository) FindSessionInfo(ctx context.Context, dto domain.FindSessionDto) (domain.SessionInfo, error) {
	if err := dto.Validate(); err != nil {
		return domain.SessionInfo{}, err
	}
	key := "sessions:" + dto.Id
	val, err := r.redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return domain.SessionInfo{}, errors.Join(ErrRedisQueryFailed, err)
	}
	if len(val) == 0 {
		return domain.SessionInfo{}, ErrSessionNotFound
	}
	return domain.SessionInfo{
		UserId:   val["user_id"],
		UserName: val["user_name"],
		UserRole: val["user_role"],
	}, nil
}

// func (r *SessionRepository) DeleteByUserId(ctx context.Context, dto domain.FindSessionDto) error {
// 	if err := dto.Validate(); err != nil {
// 		return err
// 	}

// 	query := `DELETE FROM sessions WHERE "user_id" = $1 RETURNING "session_id"`
// 	rows, err := r.pgPool.Query(ctx, query, dto.Id)
// 	if err != nil {
// 		return errors.Join(ErrPostgresQueryFailed, err)
// 	}
// 	defer rows.Close()

// 	var keys []string
// 	for rows.Next() {
// 		var sessionId string
// 		if err := rows.Scan(&sessionId); err != nil {
// 			return errors.Join(ErrPostgresQueryFailed, err)
// 		}

// 		key := "sessions:" + sessionId
// 		keys = append(keys, key)
// 	}

// 	_, err = r.redisClient.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
// 		for _, key := range keys {
// 			pipe.Del(ctx, key)
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return errors.Join(ErrRedisQueryFailed, err)
// 	}

// 	return nil
// }
