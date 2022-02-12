package entity

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       ID
	Login    string
	Password string
}

func NewUser(login string, password string) (*User, error) {
	id, err := NewID()
	if err != nil {
		return nil, err
	}
	pwd, err := generatePassword(password)
	if err != nil {
		return nil, err
	}
	u := &User{
		ID:       id,
		Login:    login,
		Password: pwd,
	}

	err = u.Validate()
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) ValidatePassword(p string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Validate() error {
	if u.Login == "" || u.Password == "" {
		return ErrInvalidEntity
	}

	return nil
}

func generatePassword(raw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
