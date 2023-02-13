package model

import (
	"errors"
	"strings"
)

type User struct {
	ID       int    `json:"-"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

var (
	ErrUserLoginEmpty    = errors.New("login empty")
	ErrUserPasswordEmpty = errors.New("password empty")
)

func (u *User) Validate() error {
	if strings.TrimSpace(u.Login) == "" {
		return ErrUserLoginEmpty
	}
	if strings.TrimSpace(u.Password) == "" {
		return ErrUserPasswordEmpty
	}
	return nil
}
