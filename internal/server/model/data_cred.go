package model

import (
	"errors"
	"strings"
	"time"
)

type DataCred struct {
	ID        int       `json:"id"`
	UserID    int       `json:"-"`
	Title     string    `json:"title"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Meta      string    `json:"meta"`
	UpdatedAt time.Time `json:"updated_at"`
}

var (
	ErrDataCredTitleEmpty    = errors.New("title empty")
	ErrDataCredUsernameEmpty = errors.New("username empty")
	ErrDataCredPasswordEmpty = errors.New("password empty")
	ErrDataCredUserIDEmpty   = errors.New("user_id empty")
)

func (d *DataCred) Validate() error {
	if strings.TrimSpace(d.Title) == "" {
		return ErrDataCredTitleEmpty
	}

	if strings.TrimSpace(d.Username) == "" {
		return ErrDataCredUsernameEmpty
	}

	if strings.TrimSpace(d.Password) == "" {
		return ErrDataCredPasswordEmpty
	}

	if d.UserID == 0 {
		return ErrDataCardUserIDEmpty
	}

	return nil
}
