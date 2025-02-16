package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kstsm/avito-shop/config"
	"log"
	"time"
)

const tokenExpiry = time.Hour * 24

func GenerateToken(userID int) (string, error) {
	secretKey := []byte(config.Config.JWT.JWTSecret)

	claims := jwt.MapClaims{
		"userID": float64(userID),
		"exp":    time.Now().Add(tokenExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	log.Printf("Генерация токена для пользователя: %d", userID)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Printf("Ошибка при подписании токена: %v\n", err)
		return "", err
	}

	log.Printf("Токен успешно сгенерирован для пользователя: %d", userID)
	return tokenString, nil
}

func ValidateToken(tokenString string) (int, error) {
	secretKey := []byte(config.Config.JWT.JWTSecret)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Ошибка: неверный метод подписи: %v\n", token.Header["alg"])
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
		log.Println("Поле exp отсутствует или неверного типа")
		return 0, errors.New("поле exp отсутствует или неверного типа")
	}
	if time.Now().Unix() > int64(expFloat) {
		log.Println("Токен истек")
		return 0, errors.New("токен истек")
	}

	userIDFloat, ok := claims["userID"].(float64)
	if !ok {
		log.Println("Поле userID отсутствует или неверного типа")
		return 0, errors.New("поле userID отсутствует или неверного типа")
	}

	userID := int(userIDFloat)
	log.Printf("Токен успешно валидирован для пользователя: %d", userID)

	return userID, nil
}
