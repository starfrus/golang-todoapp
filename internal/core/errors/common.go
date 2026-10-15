package core_errors

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrInvaligArgument = errors.New("invalid argument")
	ErrConflict = errors.New("conflict")
)
