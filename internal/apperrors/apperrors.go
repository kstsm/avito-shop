package apperrors

import "errors"

var (
	ErrUserNotFound       = errors.New("пользователь не найден")
	ErrInsufficientFunds  = errors.New("недостаточно средств")
	ErrInvalidTransfer    = errors.New("нельзя перевести монеты самому себе")
	ErrItemNotFound       = errors.New("товар не найден")
	ErrInvalidCredentials = errors.New("неверное имя пользователя или пароль")
	ErrUserCreation       = errors.New("ошибка создания пользователя")
)
