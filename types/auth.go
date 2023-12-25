package types

import "golang.org/x/crypto/bcrypt"

type AuthResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

type AuthUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func IsValidPassword(enpw, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(enpw), []byte(pw)) == nil
}
