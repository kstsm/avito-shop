package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kstsm/avito-shop/config"
	"log"
	"time"
)

func GenerateToken(userID int) (string, error) {
	secretKey := []byte(config.Config.JWT.JWTSecret)
	tokenExpiry := config.Config.JWT.TokenExpiry

	claims := jwt.MapClaims{
		"userID": float64(userID),
		"exp":    time.Now().Add(tokenExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (int, error) {
	secretKey := []byte(config.Config.JWT.JWTSecret)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неверный метод подписи")
		}
		return secretKey, nil
	})

	if err != nil {
		log.Printf("Ошибка при валидации токена: %v\n", err)
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Println("Неверный токен или поврежденные данные")
		return 0, errors.New("неверный токен")
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return 0, errors.New("поле exp отсутствует или неверного типа")
	}
	if time.Now().Unix() > int64(expFloat) {
		return 0, errors.New("токен истек")
	}

	userIDFloat, ok := claims["userID"].(float64)
	if !ok {
		return 0, errors.New("поле userID отсутствует или неверного типа")
	}

	userID := int(userIDFloat)

	return userID, nil
}
