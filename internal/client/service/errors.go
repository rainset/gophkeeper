package service

import "errors"

var (
	ErrStatusLoginExists  = errors.New("ошибка такой логин уже занят")
	ErrStatusUnauthorized = errors.New("ошибка авторизации")
	ErrServer             = errors.New("ошибка соединения с сервером")
)
