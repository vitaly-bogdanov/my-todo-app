package core_domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/vitaly-bogdanov/my-todo-app/internal/core/errors"
)

var phoneNumberPattern = regexp.MustCompile(`^\+[1-9]\d{7,14}$`)

type User struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

func NewUser(
	id int,
	version int,
	fullName string,
	phoneNumber *string,
) User {
	return User{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func NewUserUninitialized(
	fullName string,
	phoneNumber *string,
) User {
	return NewUser(
		UninitializedID,
		UninitializedVersion,
		fullName,
		phoneNumber,
	)
}

func (u *User) Validate() error {
	fullNameLength := len([]rune(u.FullName))
	if fullNameLength < 1 || fullNameLength > 100 {
		return fmt.Errorf("invalid FullName length %d: %w", fullNameLength, core_errors.BadRequest)
	}

	if u.PhoneNumber != nil && !phoneNumberPattern.MatchString(*u.PhoneNumber) {
		return fmt.Errorf("invalid PhoneNumber: %w", core_errors.BadRequest)
	}

	return nil
}

type UserPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
}

func NewUserPatch(
	fullName Nullable[string],
	phoneNumber Nullable[string],
) UserPatch {
	return UserPatch{
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func (u *UserPatch) Validate() error {
	if u.FullName.Set && u.FullName.Value == nil {
		return fmt.Errorf("FullName is required: %w", core_errors.BadRequest)
	}

	return nil
}

func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("failed to validate patch: %w", err)
	}

	tmp := *u

	if patch.FullName.Set {
		tmp.FullName = *patch.FullName.Value
	}

	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("failed to validate user: %w", err)
	}

	*u = tmp

	return nil
}
