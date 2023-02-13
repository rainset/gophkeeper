package storage

import "errors"

var (
	ErrorRowAlreadyExists = errors.New("row already exists")

	ErrorUserAlreadyExists = errors.New("user already exists")
	ErrorUserCredentials   = errors.New("wrong pair login/password")
)
