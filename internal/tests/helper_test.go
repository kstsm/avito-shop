package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kstsm/avito-shop/api/rest/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func authenticate(t *testing.T, server *httptest.Server, username, password string) (*http.Client, string) {
	authPayload := fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)
	req, err := http.NewRequest("POST", server.URL+"/api/auth", strings.NewReader(authPayload))
	require.NoError(t, err, "Ошибка при создании запроса аутентификации")
	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	require.NoError(t, err, "Ошибка при выполнении запроса аутентификации")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Ожидался статус 200 OK для запроса аутентификации")

	var authResp models.AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	require.NoError(t, err, "Ошибка при декодировании ответа аутентификации")

	require.NotEmpty(t, authResp.Token, "Токен не должен быть пустым")

	return client, "Bearer " + authResp.Token
}

func sendCoinsExpectFailure(t *testing.T, server *httptest.Server, client *http.Client, token, toUser string, amount int, expectedStatus int) *http.Response {
	sendCoinPayload := fmt.Sprintf(`{"toUser":"%s", "amount": %d}`, toUser, amount)
	req, err := http.NewRequest("POST", server.URL+"/api/sendCoin", strings.NewReader(sendCoinPayload))
	if err != nil {
		t.Fatalf("Ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса: %v", err)
	}

	if resp.StatusCode != expectedStatus {
		t.Fatalf("Ожидался статус %v, но получен %v", expectedStatus, resp.StatusCode)
	}

	return resp
}

func getUserBalance(t *testing.T, conn *pgxpool.Pool, ctx context.Context, username string) int {
	var balance int
	err := conn.QueryRow(ctx, "SELECT balance FROM users WHERE username = $1", username).Scan(&balance)
	if err != nil {
		t.Fatalf("Ошибка получения баланса пользователя '%s': %v", username, err)
	}
	return balance
}

func sendCoins(t *testing.T, server *httptest.Server, client *http.Client, token, toUser string, amount int) {
	sendCoinPayload := fmt.Sprintf(`{"toUser":"%s", "amount": %d}`, toUser, amount)
	req, err := http.NewRequest("POST", server.URL+"/api/sendCoin", strings.NewReader(sendCoinPayload))
	if err != nil {
		t.Fatalf("Ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался статус %v, но получен %v", http.StatusOK, resp.StatusCode)
	}
}
