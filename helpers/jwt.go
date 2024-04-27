package helpers

import (
	"errors"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(claims *jwt.RegisteredClaims) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")

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
