package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/starfrus/golang-todoapp/internal/core/errors"
)

type User struct {
	ID      int
	Version int

	FullName    string
	PhoneNumber *string
}

func NewUser(
	id int,
	version int,
	fullname string,
	phoneNumber *string) User {
	return User{
		ID:          id,
		Version:     version,
		FullName:    fullname,
		PhoneNumber: phoneNumber,
	}
}

func NewUserUnitialized(fullname string, phoneNumber *string) User {
	return NewUser(
		UninitializedID,
		UninitializedVersion,
		fullname,
		phoneNumber,
	)
}

func (u *User) Validate() error {
	fullnameLenght := len([]rune(u.FullName))
	if fullnameLenght < 3 || fullnameLenght > 100 {
		return fmt.Errorf(
			"invalid `FullName` len: %d %w",
			fullnameLenght, core_errors.ErrInvaligArgument,
		)
	}

	if u.PhoneNumber != nil {
		phoneNumberLenght := len([]rune(*u.PhoneNumber))
		if phoneNumberLenght < 10 || phoneNumberLenght > 15 {
			return fmt.Errorf("invalid `PhoneNumber` len: %d: %w",
				phoneNumberLenght, core_errors.ErrInvaligArgument)
		}
		re := regexp.MustCompile(`^\+[0-9]+$`)

		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf("invalid `PhoneNumber` format: %w", core_errors.ErrInvaligArgument)
		}
	}
	return nil
}

type UserPatch struct {
	FullName Nullable[string]
	PhoneNumber Nullable[string]
}

func(p *UserPatch) Validate() error{
	if p.FullName.Set && p.FullName.Value == nil{
		return fmt.Errorf(
			"'FullName' can`t be patched to NULL: %w", 
			core_errors.ErrInvaligArgument,
		)
	}
	return nil
}

func (u *User)ApplyPatch(patch UserPatch) error{
	if err := patch.Validate(); err != nil{
		return fmt.Errorf("validate user patch: %w", err)
	}

	tmp := *u

	if patch.FullName.Set{
		tmp.FullName = *patch.FullName.Value
	}

	if patch.PhoneNumber.Set{
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}

	if err := tmp.Validate(); err != nil{
		return fmt.Errorf("validate patched user: %w", err)
	}

	*u = tmp

	return nil
}