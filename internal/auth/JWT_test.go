package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/kstsm/avito-shop/config"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	config.Config.JWT.JWTSecret = "testSecretKey"

	tests := []struct {
		name   string
		userID int
	}{
		{
			name:   "Valid token generation",
			userID: 12345,
		},
		{
			name:   "Zero user ID",
			userID: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID)

			assert.NoError(t, err)
			assert.NotEmpty(t, token)

			parsedToken, parseErr := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS256 {
					t.Errorf("Неверный метод подписи: %v", token.Header["alg"])
				}
				return []byte(config.Config.JWT.JWTSecret), nil
			})

			assert.NoError(t, parseErr)
			assert.NotNil(t, parsedToken)
			assert.True(t, parsedToken.Valid)

			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			assert.True(t, ok)
			assert.Equal(t, float64(tt.userID), claims["userID"])

			expirationTime := int64(claims["exp"].(float64))
			assert.True(t, expirationTime > time.Now().Unix())
			assert.True(t, expirationTime <= time.Now().Add(tokenExpiry).Unix())
		})
	}
}

func generateTestToken(userID int, secret string, expiry time.Duration) string {
	claims := jwt.MapClaims{
		"userID": float64(userID),
		"exp":    time.Now().Add(expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func TestValidateToken(t *testing.T) {
	config.Config.JWT.JWTSecret = "testSecretKey"

	validToken := generateTestToken(123, config.Config.JWT.JWTSecret, time.Hour*24)
	expiredToken := generateTestToken(123, config.Config.JWT.JWTSecret, -time.Hour)

	invalidToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk2NDA3MjMsInVzZXJJRCI6Mzh9.vH-HnUbkupIIThFfgVUhpOmSO1QmMgjy8d326ww8ct4"

	tests := []struct {
		name       string
		token      string
		expectErr  bool
		expectedID int
		errMsg     string
	}{
		{
			name:       "Valid token",
			token:      validToken,
			expectErr:  false,
			expectedID: 123,
		},
		{
			name:      "Expired token",
			token:     expiredToken,
			expectErr: true,
			errMsg:    "token is expired",
		},
		{
			name:      "Invalid token",
			token:     invalidToken,
			expectErr: true,
			errMsg:    "token signature is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := ValidateToken(tt.token)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Equal(t, 0, userID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, userID)
			}
		})
	}
}
