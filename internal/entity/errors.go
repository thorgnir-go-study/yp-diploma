package entity

import "errors"

var (
	ErrInvalidEntity       = errors.New("invalid entity")
	ErrUserNotFound        = errors.New("user not found")
	ErrIncorrectPassword   = errors.New("incorrect password")
	ErrInsufficientBalance = errors.New("insufficient balance")
)
