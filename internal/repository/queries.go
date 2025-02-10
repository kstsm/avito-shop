package repository

const (
	GetUserByUsername = `
		SELECT id, username, password 
		FROM users 
		WHERE username=$1`

	CreateUser = `
		INSERT INTO users (username, password) 
		VALUES ($1, $2) 
		RETURNING id`
)
