package models

type UserRequestBody struct {
	Password string `json:"password"`
	Email    string `json:"email"`
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
