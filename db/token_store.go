package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ortin779/chirpy/helpers"
	"github.com/ortin779/chirpy/models"
)

func (db *DB) RefreshToken(token string) (models.RefreshTokenResponse, error) {
	dbstruct, err := db.loadDB()
	if err != nil {
		return models.RefreshTokenResponse{}, err
	}

	parsedToken, err := helpers.ValidateToken(token)

	if err != nil {
		return models.RefreshTokenResponse{}, AuthenticationError{message: err.Error()}
	}
	if !parsedToken.Valid {
		return models.RefreshTokenResponse{}, AuthenticationError{message: "invalid refresh token"}
	}
	issuer, err := parsedToken.Claims.GetIssuer()
	if err != nil {
		return models.RefreshTokenResponse{}, err
	}
	if issuer != "chirpy-refresh" {
		return models.RefreshTokenResponse{}, AuthenticationError{message: "invalid refresh token issuer"}
	}
	userId, err := parsedToken.Claims.GetSubject()
	if err != nil {
		return models.RefreshTokenResponse{}, errors.New("invalid refresh token claims")
	}

	rToken, ok := dbstruct.RefreshToken[parsedToken.Raw]
	if !ok {
		return models.RefreshTokenResponse{}, AuthenticationError{message: "invalid refresh token"}
	}

	if rToken.HasRevoked {
		return models.RefreshTokenResponse{}, AuthenticationError{message: "refresh token revoked"}
	}

	accessTokenClaims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExpiry)),
		Issuer:    "chirpy-access",
		Subject:   userId,
	}

	accessToken, err := helpers.CreateToken(accessTokenClaims)
	if err != nil {
		return models.RefreshTokenResponse{}, errors.New("error while signing the token")
	}
	return models.RefreshTokenResponse{Token: accessToken}, nil
}

func (db *DB) RevokeToken(token string) error {
	dbstruct, err := db.loadDB()
	if err != nil {
		return err
	}

	parsedToken, err := helpers.ValidateToken(token)

	if err != nil {
		return AuthenticationError{message: err.Error()}
	}
	if !parsedToken.Valid {
		return AuthenticationError{message: "invalid refresh token"}
	}
	issuer, err := parsedToken.Claims.GetIssuer()
	if err != nil {
		return err
	}
	if issuer != "chirpy-refresh" {
		return AuthenticationError{message: "invalid refresh token issuer"}
	}

	rToken, ok := dbstruct.RefreshToken[parsedToken.Raw]
	if !ok {
		return AuthenticationError{message: "invalid refresh token"}
	}

	if rToken.HasRevoked {
		fmt.Println("Revoked state")
		return AuthenticationError{message: "token has been revoked"}
	}

	dbstruct.RefreshToken[parsedToken.Raw] = models.RefreshToken{
		Id:         rToken.Id,
		HasRevoked: true,
	}

	return db.writeDB(dbstruct)

}
