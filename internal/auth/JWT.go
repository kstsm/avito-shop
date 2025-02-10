package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gookit/slog"
	"github.com/kstsm/avito-shop/config"
	"time"
)

//var secretKey = []byte("superstructure")

func GenerateToken(username string) (string, error) {
	// Получаем secretKey и tokenExpiry из конфигурации
	secretKey := []byte(config.Config.JWTSecret.JWTSecret) // Преобразуем строку в байты для подписи
	tokenExpiry := config.Config.JWTSecret.TokenExpiry     // Время жизни токена (уже как time.Duration)

	// Создаем claims для токена
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(tokenExpiry).Unix(), // Используем tokenExpiry из конфигурации
	}

	// Генерируем новый токен с указанием метода подписи
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Возвращаем подписанный токен
	return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (string, error) {
	// Получаем секретный ключ из конфигурации
	secretKey := []byte(config.Config.JWTSecret.JWTSecret)

	// Парсим токен, проверяя подпись
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверка метода подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неверный метод подписи")
		}
		return secretKey, nil
	})
	if err != nil {
		slog.Errorf("Ошибка при валидации токена: %v", err)
		return "", err
	}

	// Проверка валидности токена (включает проверку exp)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Проверяем срок действия токена
		exp := claims["exp"].(float64)
		if time.Now().Unix() > int64(exp) {
			slog.Warnf("Токен истек в %v", time.Unix(int64(exp), 0))
			return "", errors.New("токен истек")
		}
		return claims["username"].(string), nil
	}

	slog.Warn("Неверный токен или невалидные данные в токене")
	return "", errors.New("неверный токен")
}
