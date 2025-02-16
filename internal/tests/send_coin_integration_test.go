package tests

import (
	"context"
	"encoding/json"
	"github.com/kstsm/avito-shop/api/rest/models"
	"github.com/kstsm/avito-shop/database"
	"github.com/kstsm/avito-shop/internal/repository"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestSuccessfulTransfer(t *testing.T) {
	server, ctx, pool := SetupTestServer(t)
	defer server.Close()
	t.Cleanup(func() {
		pool.Close()
	})

	_, err := pool.Exec(ctx, `UPDATE users SET balance = 100 WHERE username = $1`, "senderuser")
	if err != nil {
		t.Fatalf("Ошибка обновления баланса: %v", err)
	}

	client, token := authenticate(t, server, "senderuser", "password")
	authenticate(t, server, "receiveruser", "password")

	senderBalanceBefore := getUserBalance(t, pool, ctx, "senderuser")
	receiverBalanceBefore := getUserBalance(t, pool, ctx, "receiveruser")

	sendCoins(t, server, client, token, "receiveruser", 100)

	senderBalanceAfter := getUserBalance(t, pool, ctx, "senderuser")
	receiverBalanceAfter := getUserBalance(t, pool, ctx, "receiveruser")

	if senderBalanceAfter != senderBalanceBefore-100 {
		t.Fatalf("Баланс отправителя некорректен: ожидался %d, но получен %d", senderBalanceBefore-100, senderBalanceAfter)
	}

	if receiverBalanceAfter != receiverBalanceBefore+100 {
		t.Fatalf("Баланс получателя некорректен: ожидался %d, но получен %d", receiverBalanceBefore+100, receiverBalanceAfter)
	}
}

func TestTransferToSelf(t *testing.T) {
	server, _, _ := SetupTestServer(t)
	defer server.Close()

	client, token := authenticate(t, server, "senderuser", "password")
	resp := sendCoinsExpectFailure(t, server, client, token, "senderuser", 1, http.StatusBadRequest)

	var errorResp models.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&errorResp)
	if err != nil {
		t.Fatalf("Ошибка декодирования ответа: %v", err)
	}
}

func TestInsufficientBalance(t *testing.T) {
	server, ctx, conn := SetupTestServer(t)
	defer server.Close()
	t.Cleanup(func() {
		conn.Close()
	})

	client, token := authenticate(t, server, "senderuser", "password")

	_, err := conn.Exec(ctx, "UPDATE users SET balance = 50 WHERE username = $1", "senderuser")
	if err != nil {
		t.Fatalf("Ошибка при установке баланса отправителя: %v", err)
	}

	senderBalanceBefore := getUserBalance(t, conn, ctx, "senderuser")
	if senderBalanceBefore != 50 {
		t.Fatalf("Ожидался баланс отправителя 50, но получили: %d", senderBalanceBefore)
	}

	resp := sendCoinsExpectFailure(t, server, client, token, "receiveruser", 100, http.StatusBadRequest)

	var errorResp models.ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errorResp)
	if err != nil {
		t.Fatalf("Ошибка декодирования ответа: %v", err)
	}

	if errorResp.Errors != "Недостаточно средств" {
		t.Fatalf("Ожидалось сообщение об ошибке о недостаточном балансе, но получили: %v", errorResp.Errors)
	}
}

func TestTransferZeroCoins(t *testing.T) {
	server, ctx, conn := SetupTestServer(t)
	defer server.Close()
	t.Cleanup(func() {
		conn.Close()
	})
	_, err := conn.Exec(ctx, "UPDATE users SET balance = 0 WHERE username = $1", "senderuser")

	client, token := authenticate(t, server, "senderuser", "password")
	resp := sendCoinsExpectFailure(t, server, client, token, "receiveruser", 0, http.StatusBadRequest)

	var errorResp models.ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errorResp)
	if err != nil {
		t.Fatalf("Ошибка декодирования ответа: %v", err)
	}
}

func TestTransferToNonExistentUser(t *testing.T) {
	server, ctx, conn := SetupTestServer(t)
	defer server.Close()
	t.Cleanup(func() {
		conn.Close()
	})
	_, err := conn.Exec(ctx, "UPDATE users SET balance = 100 WHERE username = $1", "senderuser")

	client, token := authenticate(t, server, "senderuser", "password")

	resp := sendCoinsExpectFailure(t, server, client, token, "nonexistentuser", 1, http.StatusNotFound)

	conn.Exec(ctx, "UPDATE users SET balance = 0 WHERE username = $1", "senderuser")

	var errorResp models.ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errorResp)
	if err != nil {
		t.Fatalf("Ошибка декодирования ответа: %v", err)
	}
}

func TestSendCoin_500(t *testing.T) {
	ctx := context.Background()
	pool, err := database.InitTestPostgres(ctx)
	require.NoError(t, err)
	defer pool.Close()

	repo := repository.NewRepository(pool)

	pool.Close()
	err = repo.SendCoins(ctx, 1, 2, "username")

	require.Error(t, err)
	require.Contains(t, err.Error(), "r.conn.Begin", "Ожидаемая ошибка, связанная с проблемой соединения с БД")
}
