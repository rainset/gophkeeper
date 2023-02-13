package model

import "time"

type User struct {
	ID       int    `json:"-"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ConfigStorage struct {
	Login        string
	Password     string
	AccessToken  string
	RefreshToken string
}

type DataCard struct {
	LocalID    int       `storm:"id,increment"`
	ExternalID int       `json:"id" storm:"unique"`
	Title      string    `json:"title"`
	Number     string    `json:"number"`
	Date       string    `json:"date"`
	Cvv        string    `json:"cvv"`
	Meta       string    `json:"meta"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DataCred struct {
	LocalID    int       `storm:"id,increment"`
	ExternalID int       `json:"id" storm:"unique"`
	Title      string    `json:"title"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Meta       string    `json:"meta"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DataText struct {
	LocalID    int       `storm:"id,increment"`
	ExternalID int       `json:"id" storm:"unique"`
	Title      string    `json:"title"`
	Text       string    `json:"text"`
	Meta       string    `json:"meta"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DataFile struct {
	LocalID    int       `storm:"id,increment"`
	ExternalID int       `json:"id" storm:"unique"`
	Title      string    `json:"title"`
	Filename   string    `json:"filename"`
	Path       string    `json:"path"`
	Ext        string    `json:"-"`
	Meta       string    `json:"meta"`
	UpdatedAt  time.Time `json:"updated_at"`
}
