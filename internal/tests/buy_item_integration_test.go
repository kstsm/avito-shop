package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kstsm/avito-shop/api/rest/models"
	"github.com/kstsm/avito-shop/database"
	"github.com/kstsm/avito-shop/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestBuyItem(t *testing.T) {
	server, ctx, conn := SetupTestServer(t)
	defer server.Close()

	t.Cleanup(func() {
		if err := conn.Close(ctx); err != nil {
			t.Errorf("Ошибка при закрытии подключения к базе данных: %v", err)
		}
	})

	client, token := authenticate(t, server, "buyer", "password")

	_, err := conn.Exec(ctx, `UPDATE users SET balance = 500 WHERE username = $1`, "buyer")
	require.NoError(t, err, "Ошибка при обновлении баланса пользователя")

	req, err := http.NewRequest("GET", server.URL+"/api/info", nil)
	require.NoError(t, err, "Ошибка при создании запроса")
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	require.NoError(t, err, "Ошибка при выполнении запроса")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Ожидался статус 200 OK")

	var infoResp models.InfoResponse
	err = json.NewDecoder(resp.Body).Decode(&infoResp)
	require.NoError(t, err, "Ошибка при декодировании ответа")

	item := "t-shirt"
	req, err = http.NewRequest("GET", fmt.Sprintf("%s/api/buy/%s", server.URL, item), nil)
	require.NoError(t, err, "Ошибка при создании запроса на покупку")
	req.Header.Set("Authorization", token)

	resp, err = client.Do(req)
	require.NoError(t, err, "Ошибка при выполнении запроса на покупку")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Ожидался статус 200 OK после покупки")

	req, err = http.NewRequest("GET", server.URL+"/api/info", nil)
	require.NoError(t, err, "Ошибка при создании запроса для проверки баланса")
	req.Header.Set("Authorization", token)

	resp, err = client.Do(req)
	require.NoError(t, err, "Ошибка при выполнении запроса для проверки баланса")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Ожидался статус 200 OK после покупки")

	err = json.NewDecoder(resp.Body).Decode(&infoResp)
	require.NoError(t, err, "Ошибка при декодировании ответа после покупки")
}

func TestBuyItemInsufficientFunds(t *testing.T) {
	server, ctx, conn := SetupTestServer(t)
	defer server.Close()

	t.Cleanup(func() {
		if err := conn.Close(ctx); err != nil {
			t.Errorf("Ошибка при закрытии подключения к базе данных: %v", err)
		}
	})

	client, token := authenticate(t, server, "testuser", "password")

	_, err := conn.Exec(ctx, `UPDATE users SET balance = 0 WHERE username = $1`, "testuser")
	require.NoError(t, err, "Ошибка при обновлении баланса пользователя")

	req, err := http.NewRequest("GET", server.URL+"/api/buy/t-shirt", nil)
	require.NoError(t, err, "Ошибка при создании запроса на покупку")
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	require.NoError(t, err, "Ошибка при выполнении запроса на покупку")
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode, "Ожидался статус 400 Bad Request")
}

func TestBuyItemNotFound(t *testing.T) {
	server, ctx, conn := SetupTestServer(t)
	defer server.Close()

	t.Cleanup(func() {
		if err := conn.Close(ctx); err != nil {
			t.Errorf("Ошибка при закрытии подключения к базе данных: %v", err)
		}
	})

	client, token := authenticate(t, server, "testuser", "password")

	req, err := http.NewRequest("GET", server.URL+"/api/buy/nonexistent-item", nil)
	require.NoError(t, err, "Ошибка при создании запроса на покупку")
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	require.NoError(t, err, "Ошибка при выполнении запроса на покупку")
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode, "Ожидался статус 404 Not Found")

	var errResp models.ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errResp)
	require.NoError(t, err, "Ошибка при декодировании ошибки")
	assert.Equal(t, "Товар не найден", errResp.Errors, "Ошибка при сообщении")
}

func TestBuyItemUnauthorized(t *testing.T) {
	server, ctx, conn := SetupTestServer(t)
	defer server.Close()

	t.Cleanup(func() {
		if err := conn.Close(ctx); err != nil {
			t.Errorf("Ошибка при закрытии подключения к базе данных: %v", err)
		}
	})
	client := http.DefaultClient
	req, err := http.NewRequest("GET", server.URL+"/api/buy/t-shirt", nil)
	require.NoError(t, err, "Ошибка при создании запроса")

	resp, err := client.Do(req)
	require.NoError(t, err, "Ошибка при выполнении запроса")

	require.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Ожидался статус 401 Unauthorized")
}

func TestBuyItem_500(t *testing.T) {
	ctx := context.Background()
	conn, err := database.InitTestPostgres(ctx)
	require.NoError(t, err)
	defer conn.Close(ctx)

	repo := repository.NewRepository(conn)

	conn.Close(ctx)

	err = repo.BuyItem(ctx, 1, "test_item")
	require.Error(t, err)
	require.Contains(t, err.Error(), "r.conn.Begin", "Ожидаемая ошибка, связанная с проблемой соединения с БД")
}
