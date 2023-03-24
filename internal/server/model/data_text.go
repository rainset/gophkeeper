package model

import (
	"errors"
	"strings"
	"time"
)

type DataText struct {
	ID        int       `json:"id"`
	UserID    int       `json:"-"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	Meta      string    `json:"meta"`
	UpdatedAt time.Time `json:"updated_at"`
}

var (
	ErrDataTextTitleEmpty  = errors.New("title empty")
	ErrDataTextEmpty       = errors.New("text empty")
	ErrDataTextUserIDEmpty = errors.New("user id empty")
)

func (d *DataText) Validate() error {
	if strings.TrimSpace(d.Title) == "" {
		return ErrDataTextTitleEmpty
	}

	if strings.TrimSpace(d.Text) == "" {
		return ErrDataTextEmpty
	}

	if d.UserID == 0 {
		return ErrDataTextUserIDEmpty
	}

	return nil
}
