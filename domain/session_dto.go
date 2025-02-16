package domain

import (
	"errors"
)

type CreateSessionDto struct {
	UserId      string
	SessionRole string
}

type FindSessionDto struct {
	Id string
}

type FindSessionWithRoleDto struct {
	Id          string
	SessionRole string
}

func (dto CreateSessionDto) Validate() error {
	var validationErrors []error

	if err := ValidateUuid(dto.UserId); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := ValidateUserRole(dto.SessionRole); err != nil {
		validationErrors = append(validationErrors, err)
	}

	if len(validationErrors) > 0 {
		return errors.Join(
			ErrValidationError,
			errors.Join(validationErrors...),
		)
	}

	return nil
}

func (dto FindSessionDto) Validate() error {
	var validationErrors []error

	if err := ValidateUuid(dto.Id); err != nil {
		validationErrors = append(validationErrors, err)
	}

	if len(validationErrors) > 0 {
		return errors.Join(
			ErrValidationError,
			errors.Join(validationErrors...),
		)
	}

	return nil
}

func (dto FindSessionWithRoleDto) Validate() error {
	var validationErrors []error

	if err := ValidateUuid(dto.Id); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := ValidateUserRole(dto.SessionRole); err != nil {
		validationErrors = append(validationErrors, err)
	}

	if len(validationErrors) > 0 {
		return errors.Join(
			ErrValidationError,
			errors.Join(validationErrors...),
		)
	}

	return nil
}
