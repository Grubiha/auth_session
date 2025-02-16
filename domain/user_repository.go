package domain

import "context"

type UserRepository interface {
	Create(ctx context.Context, dto CreateUserDto) (string, error)
	Find(ctx context.Context, dto FindUserDto) (User, error)
	FindByPhone(ctx context.Context, dto FindUserByPhoneDto) (User, error)
	Update(ctx context.Context, dto UpdateUserDto) error
	Delete(ctx context.Context, dto FindUserDto) error
}
