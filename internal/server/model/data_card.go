package model

import (
	"errors"
	"strings"
	"time"
)

type DataCard struct {
	ID        int       `json:"id"`
	UserID    int       `json:"-"`
	Title     string    `json:"title"`
	Number    string    `json:"number"`
	Date      string    `json:"date"`
	Cvv       string    `json:"cvv"`
	Meta      string    `json:"meta"`
	UpdatedAt time.Time `json:"updated_at"`
}

var (
	ErrDataCardTitleEmpty  = errors.New("title empty")
	ErrDataCardUserIDEmpty = errors.New("user id empty")
	ErrDataCardNumberEmpty = errors.New("number empty")
	ErrDataCardDateEmpty   = errors.New("date empty")
	ErrDataCardCvvEmpty    = errors.New("cvv empty")
)

func (d *DataCard) Validate() error {
	if strings.TrimSpace(d.Title) == "" {
		return ErrDataCardTitleEmpty
	}
	if d.UserID == 0 {
		return ErrDataCardUserIDEmpty
	}
	if strings.TrimSpace(d.Number) == "" {
		return ErrDataCardNumberEmpty
	}
	if strings.TrimSpace(d.Date) == "" {
		return ErrDataCardDateEmpty
	}
	if strings.TrimSpace(d.Cvv) == "" {
		return ErrDataCardCvvEmpty
	}
	return nil
}
