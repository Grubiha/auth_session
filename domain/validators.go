package domain

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"
)

func ValidateUuid(stringID string) error {
	id, err := uuid.Parse(stringID)
	if err != nil {
		return fmt.Errorf(`%w expected valid uuid`, ErrInvalidUuid)
	}
	if id == uuid.Nil {
		return fmt.Errorf(`%w expected non-nil "id"`, ErrInvalidUuid)
	}
	return nil
}

func ValidateUserName(name string) error {
	if name == "" {
		return fmt.Errorf(`%w expected non-empty "name"`, ErrInvalidUserName)
	}
	nameRegex := regexp.MustCompile(`^[a-zA-Zа-яА-Я\s]{1,100}$`)
	if !nameRegex.MatchString(name) {
		return fmt.Errorf(`%w expected ^[a-zA-Zа-яА-Я\s]{1,100}$`, ErrInvalidUserName)
	}
	return nil
}

func ValidateUserPhone(phone string) error {
	phoneRegex := regexp.MustCompile(`^\+7\d{10}$`)
	if !phoneRegex.MatchString(phone) {
		return fmt.Errorf(`%w expected ^\+7\d{10}$`, ErrInvalidUserPhone)
	}
	return nil
}

func ValidateUserRole(role string) error {

	if !ValidUserRoles[role] {
		return fmt.Errorf(`%w expected one of: "user", "admin", "manager"`, ErrInvalidUserRole)
	}
	return nil
}
