package service

import (
	"github.com/kstsm/avito-shop/internal/auth"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestAuthenticate_NewUser_HashingAndToken(t *testing.T) {

	password := "securepassword"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	assert.NoError(t, err)

	token, err := auth.GenerateToken(1)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
