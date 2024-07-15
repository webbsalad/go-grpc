package storage

import "errors"

var (
	ErrUserExist    = errors.New("пользователь уже зарегестрирован")
	ErrUserNotFound = errors.New("пользователь не найден")
	ErrAppNotFound  = errors.New("приложение не найдено")
)
