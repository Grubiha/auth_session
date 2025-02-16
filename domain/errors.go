package domain

import "errors"

var (
	ErrValidationError  = errors.New("validation error")
	ErrInvalidUuid      = errors.New("invalid uuid")
	ErrInvalidUserName  = errors.New("invalid user name")
	ErrInvalidUserPhone = errors.New("invalid user phone")
	ErrInvalidUserRole  = errors.New("invalid user role")
)
