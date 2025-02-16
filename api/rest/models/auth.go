package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}
type AuthRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30,excludesrune= ,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=20,excludesrune= "`
}

type AuthResponse struct {
	Token string `json:"token"`
}
