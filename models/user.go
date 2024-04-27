package models

type UserRequestBody struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserLoginResponse struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
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
