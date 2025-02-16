package repository

const (
	QueryGetUserByUsername = `
		SELECT id, username, password 
		FROM users 
		WHERE username=$1`

	QueryCreateUser = `
		INSERT INTO users (username, password) 
		VALUES ($1, $2) 
		RETURNING id`

	QueryGetItem = `
		SELECT id, price 
		FROM items 
		WHERE name = $1`

	QueryUpdateUserBalance = `
		UPDATE users 
		SET balance = balance - $1 
		WHERE id = $2`

	QueryUpdateInventory = `
		INSERT INTO inventory (user_id, item_id, quantity) 
		VALUES ($1, $2, 1) 
		ON CONFLICT (user_id, item_id) DO UPDATE 
		SET quantity = inventory.quantity + 1`

	QueryInsertTransaction = `
		INSERT INTO transactions (from_user, to_user, amount) 
		VALUES ($1, $2, $3);
	`
	QueryTransferCoins = `
		WITH updated AS (
    UPDATE users
    SET balance = CASE
                    WHEN id = $1 THEN balance - $2 
                    WHEN id = $3 THEN balance + $2
                  END
    WHERE id IN ($1, $3)
    RETURNING id
)
-- Проверим, что обновились две строки (одна для отправителя и одна для получателя)
SELECT COUNT(*)
FROM updated
WHERE id IN ($1, $3);`

	QueryGetUserInfo = `
SELECT json_build_object(
  'coins', u.balance,
  'inventory', (
    SELECT COALESCE(json_agg(json_build_object('type', i.name, 'quantity', inv.quantity)), '[]'::json)
    FROM inventory inv
    JOIN items i ON inv.item_id = i.id
    WHERE inv.user_id = u.id
  ),
  'coinHistory', json_build_object(
    'received', (
      SELECT COALESCE(json_agg(json_build_object('fromUser', u2.username, 'amount', t.amount)), '[]'::json)
      FROM transactions t
      JOIN users u2 ON t.from_user = u2.id
      WHERE t.to_user = u.id
    ),
    'sent', (
      SELECT COALESCE(json_agg(json_build_object('toUser', u2.username, 'amount', t.amount)), '[]'::json)
      FROM transactions t
      JOIN users u2 ON t.to_user = u2.id
      WHERE t.from_user = u.id
    )
  )
) 
FROM users u 
WHERE u.id = $1;
`
)
