package db

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ortin779/chirpy/helpers"
	"github.com/ortin779/chirpy/models"
	"golang.org/x/crypto/bcrypt"
)

func (db *DB) CreateUser(userBody models.UserRequestBody) (models.UserResponse, error) {
	dbstruct, err := db.loadDB()
	if err != nil {
		return models.UserResponse{}, err
	}

	existingUsr := findUser(userBody.Email, dbstruct.Users)
	if existingUsr != nil {
		return models.UserResponse{}, fmt.Errorf("user already exist with given email")
	}

	nextIndex := 1

	if len(dbstruct.Users) > 0 {
		keys := getSortedKeys(dbstruct.Users)
		nextIndex = keys[0] + 1
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userBody.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.UserResponse{}, err
	}
	newUser := models.User{
		Id:       nextIndex,
		Email:    userBody.Email,
		Password: string(hashedPassword),
	}
	dbstruct.Users[nextIndex] = newUser
	err = db.writeDB(dbstruct)
	if err != nil {
		return models.UserResponse{}, err
	}
	return models.UserResponse{
		Id:    newUser.Id,
		Email: newUser.Email,
	}, nil
}

func (db *DB) UpdateUser(userBody models.UserRequestBody, userId string) (models.UserResponse, error) {
	dbstruct, err := db.loadDB()
	if err != nil {
		return models.UserResponse{}, err
	}

	parsedId, err := strconv.Atoi(userId)
	if err != nil {
		return models.UserResponse{}, fmt.Errorf("invalid user id")
	}

	existingUsr, ok := dbstruct.Users[parsedId]
	if !ok {
		return models.UserResponse{}, NotFoundError{}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userBody.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.UserResponse{}, err
	}
	updatedUser := models.User{
		Id:       existingUsr.Id,
		Email:    userBody.Email,
		Password: string(hashedPassword),
	}
	dbstruct.Users[parsedId] = updatedUser
	err = db.writeDB(dbstruct)
	if err != nil {
		return models.UserResponse{}, err
	}
	return models.UserResponse{
		Id:    updatedUser.Id,
		Email: updatedUser.Email,
	}, nil
}

func (db *DB) LoginUser(userBody models.UserRequestBody) (models.UserLoginResponse, error) {
	dbstruct, err := db.loadDB()
	if err != nil {
		return models.UserLoginResponse{}, err
	}

	user := findUser(userBody.Email, dbstruct.Users)
	if user == nil {
		return models.UserLoginResponse{}, AuthenticationError{message: fmt.Sprintf("no user with given email %s", userBody.Email)}
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userBody.Password))
	if err != nil {
		return models.UserLoginResponse{}, AuthenticationError{message: fmt.Sprintf("invalid password for user with email %s", userBody.Email)}
	}

	accessTokenClaims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExpiry)),
		Issuer:    "chirpy-access",
		Subject:   strconv.Itoa(user.Id),
	}

	refreshTokenClaims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenExpiry)),
		Issuer:    "chirpy-refresh",
		Subject:   strconv.Itoa(user.Id),
	}

	accessToken, err := helpers.CreateToken(accessTokenClaims)
	if err != nil {
		return models.UserLoginResponse{}, errors.New("error while signing the token")
	}
	refreshToken, err := helpers.CreateToken(refreshTokenClaims)
	if err != nil {
		return models.UserLoginResponse{}, errors.New("error while signing the token")
	}
	dbstruct.RefreshToken[refreshToken] = models.RefreshToken{
		Id:         refreshToken,
		HasRevoked: false,
	}
	err = db.writeDB(dbstruct)
	if err != nil {
		return models.UserLoginResponse{}, err
	}
	return models.UserLoginResponse{
		Id:           user.Id,
		Email:        user.Email,
		Token:        accessToken,
		RefreshToken: refreshToken,
		IsChirpyRed:  user.IsChirpyRed,
	}, nil
}

func (db *DB) MarkUserAsRedChirp(userId int) error {
	dbstruct, err := db.loadDB()
	if err != nil {
		return err
	}

	existingUsr, ok := dbstruct.Users[userId]
	if !ok {
		return NotFoundError{}
	}

	updatedUser := models.User{
		Id:          existingUsr.Id,
		Email:       existingUsr.Email,
		Password:    existingUsr.Password,
		IsChirpyRed: true,
	}

	dbstruct.Users[userId] = updatedUser

	err = db.writeDB(dbstruct)
	if err != nil {
		return err
	}
	return nil
}
