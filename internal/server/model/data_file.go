package model

import (
	"errors"
	"strings"
	"time"
)

type DataFile struct {
	ID        int       `json:"id"`
	UserID    int       `json:"-"`
	Title     string    `json:"title"`
	Filename  string    `json:"filename"`
	Path      string    `json:"path"`
	Meta      string    `json:"meta"`
	UpdatedAt time.Time `json:"updated_at"`
}

var (
	ErrDataFileTitleEmpty  = errors.New("title empty")
	ErrDataFilePathEmpty   = errors.New("file path empty")
	ErrDataFileUserIDEmpty = errors.New("user id empty")
)

func (d *DataFile) Validate() error {
	if strings.TrimSpace(d.Path) == "" {
		return ErrDataFilePathEmpty
	}

	if strings.TrimSpace(d.Title) == "" {
		return ErrDataFileTitleEmpty
	}

	if d.UserID == 0 {
		return ErrDataFileUserIDEmpty
	}

	return nil
}
