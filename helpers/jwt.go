package helpers

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	. "github.com/ortin779/chirpy/models"
)

func CreateToken(userBody UserRequestBody, sub string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	expiryTime := 50_400
	if userBody.ExpiresInSec != 0 && userBody.ExpiresInSec <= expiryTime {
		expiryTime = userBody.ExpiresInSec
	}
	claims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(expiryTime)).UTC()),
		Issuer:    "chirpy",
		Subject:   sub,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecret))
}

func ValidateToken(token string) (*jwt.Token, error) {
	jwtSecret := os.Getenv("JWT_SECRET")

	tokenParts := strings.Split(token, " ")

	if len(tokenParts) != 2 {
		return nil, errors.New("invalid token")
	}

	return jwt.ParseWithClaims(tokenParts[1], &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
}
