package models

type RefreshTokenResponse struct {
	Token string `json:"token"`
}

type RefreshToken struct {
	Id         string `json:"id"`
	HasRevoked bool   `json:"hasRevoked"`
}
