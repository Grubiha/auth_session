package repos

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Grubiha/auth_session/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) domain.UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (r *UserRepository) Create(ctx context.Context, dto domain.CreateUserDto) (string, error) {
	if err := dto.Validate(); err != nil {
		return "", err
	}

	var query string
	var args []interface{}

	// Формируем SQL-запрос в зависимости от наличия роли
	if dto.Role == nil {
		query = `INSERT INTO users ("user_name", "user_phone") VALUES ($1, $2) RETURNING "user_id"`
		args = []interface{}{dto.Name, dto.Phone}
	} else {
		query = `INSERT INTO users ("user_name", "user_phone", "user_role") VALUES ($1, $2, $3) RETURNING "user_id"`
		args = []interface{}{dto.Name, dto.Phone, dto.Role}
	}

	var id string
	err := r.pool.QueryRow(ctx, query, args...).Scan(&id)

	// Обрабатываем ошибку нарушения уникальности (код ошибки 23505)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			return "", ErrUniqueViolation
		}
		return "", errors.Join(ErrPostgresQueryFailed, err) // ErrQueryFailed
	}

	return id, nil
}

func (r *UserRepository) Delete(ctx context.Context, dto domain.FindUserDto) error {
	if err := dto.Validate(); err != nil {
		return err
	}

	query := `DELETE FROM users WHERE "user_id" = $1`

	result, err := r.pool.Exec(ctx, query, dto.Id)
	if err != nil {
		return errors.Join(ErrPostgresQueryFailed, err)
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil

}

func (r *UserRepository) Find(ctx context.Context, dto domain.FindUserDto) (domain.User, error) {
	if err := dto.Validate(); err != nil {
		return domain.User{}, err
	}

	query := `SELECT "user_id", "user_name", "user_phone", "user_role" FROM users WHERE "user_id" = $1`

	var user domain.User
	err := r.pool.QueryRow(ctx, query, dto.Id).Scan(&user.Id, &user.Name, &user.Phone, &user.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, ErrUserNotFound
		}
		return domain.User{}, errors.Join(ErrPostgresQueryFailed, err)
	}

	return user, nil
}

func (r *UserRepository) FindByPhone(ctx context.Context, dto domain.FindUserByPhoneDto) (domain.User, error) {
	if err := dto.Validate(); err != nil {
		return domain.User{}, err
	}

	query := `SELECT "user_id", "user_name", "user_phone", "user_role" FROM users WHERE "user_phone" = $1`

	var user domain.User
	err := r.pool.QueryRow(ctx, query, dto.Phone).Scan(&user.Id, &user.Name, &user.Phone, &user.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, ErrUserNotFound
		}
		return domain.User{}, errors.Join(ErrPostgresQueryFailed, err)
	}

	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, dto domain.UpdateUserDto) error {
	if err := dto.Validate(); err != nil {
		return err
	}

	// Подготавливаем запрос
	var args []interface{}
	var queryParts []string
	var queryColumns []string

	if dto.Name != nil {
		queryColumns = append(queryColumns, `"user_name"`)
		args = append(args, *dto.Name)
		queryParts = append(queryParts, fmt.Sprintf("$%d", len(args)))
	}

	if dto.Phone != nil {
		queryColumns = append(queryColumns, `"user_phone"`)
		args = append(args, *dto.Phone)
		queryParts = append(queryParts, fmt.Sprintf("$%d", len(args)))
	}

	if dto.Role != nil {
		queryColumns = append(queryColumns, `"user_role"`)
		args = append(args, *dto.Role)
		queryParts = append(queryParts, fmt.Sprintf("$%d", len(args)))
	}

	// Генерируем часть SET для SQL-запроса
	querySet := make([]string, len(queryParts))
	for i, part := range queryParts {
		querySet[i] = fmt.Sprintf(`%s = %s`, queryColumns[i], part)
	}

	// Добавляем ID как последний аргумент
	args = append(args, dto.Id)

	query := fmt.Sprintf(
		`UPDATE users SET %s WHERE "user_id" = $%d`,
		strings.Join(querySet, ", "),
		len(args),
	)

	result, err := r.pool.Exec(ctx, query, args...)

	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			return ErrUniqueViolation
		}
		return errors.Join(ErrPostgresQueryFailed, err)
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil

}
