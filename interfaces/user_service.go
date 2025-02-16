package interfaces

import (
	"context"

	"github.com/Grubiha/auth_session/domain"
)

type UserService interface {
	Create(ctx context.Context, dto domain.CreateUserDto) (string, error)
	Find(ctx context.Context, dto domain.FindUserDto) (domain.User, error)
	FindByPhone(ctx context.Context, dto domain.FindUserByPhoneDto) (domain.User, error)
	Update(ctx context.Context, dto domain.UpdateUserDto) error
	Delete(ctx context.Context, dto domain.FindUserDto) error
}
