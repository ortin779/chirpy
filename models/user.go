package models

type UserRequestBody struct {
	Password string `json:"password"`
	Email    string `json:"email"`

	//OptionalFields
	ExpiresInSec int `json:"expires_in_seconds"`
}

type UserLoginResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type UserResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
