package domain

import (
	"errors"
	"strings"
)

type CreateUserDto struct {
	Name  string
	Phone string
	Role  *string
}

type FindUserDto struct {
	Id string
}

type FindUserByPhoneDto struct {
	Phone string
}

type UpdateUserDto struct {
	Id    string
	Name  *string
	Phone *string
	Role  *string
}

func (dto CreateUserDto) Validate() error {
	var validationErrors []error
	// Transform
	dto.Name = strings.Join(strings.Fields(dto.Name), " ")
	// Validate
	if err := ValidateUserName(dto.Name); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if err := ValidateUserPhone(dto.Phone); err != nil {
		validationErrors = append(validationErrors, err)
	}
	if dto.Role != nil {
		if err := ValidateUserRole(*dto.Role); err != nil {
			validationErrors = append(validationErrors, err)
		}
	}
	if len(validationErrors) > 0 {
		return errors.Join(
			ErrValidationError,
			errors.Join(validationErrors...),
		)
	}
	return nil
}

func (dto FindUserDto) Validate() error {
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

func (dto FindUserByPhoneDto) Validate() error {
	var validationErrors []error

	if err := ValidateUserPhone(dto.Phone); err != nil {
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

func (dto UpdateUserDto) Validate() error {
	var validationErrors []error

	if err := ValidateUuid(dto.Id); err != nil {
		validationErrors = append(validationErrors, err)
	}

	if dto.Name != nil {
		if err := ValidateUserName(*dto.Name); err != nil {
			validationErrors = append(validationErrors, err)
		}
	}
	if dto.Phone != nil {
		if err := ValidateUserPhone(*dto.Phone); err != nil {
			validationErrors = append(validationErrors, err)
		}
	}
	if dto.Role != nil {
		if err := ValidateUserRole(*dto.Role); err != nil {
			validationErrors = append(validationErrors, err)
		}
	}

	if len(validationErrors) > 0 {
		return errors.Join(
			ErrValidationError,
			errors.Join(validationErrors...),
		)
	}

	return nil
}
