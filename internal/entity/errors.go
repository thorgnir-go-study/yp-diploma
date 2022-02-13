package entity

import "errors"

var (
	ErrUserAlreadyExists            = errors.New("user already exists")
	ErrInvalidEntity                = errors.New("invalid entity")
	ErrUserNotFound                 = errors.New("user not found")
	ErrIncorrectPassword            = errors.New("incorrect password")
	ErrInsufficientBalance          = errors.New("insufficient balance")
	ErrOrderAlreadyRegistered       = errors.New("order already registered")
	ErrOrderRegisteredByAnotherUser = errors.New("order already registered by another user")
)
