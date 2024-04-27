package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ortin779/chirpy/helpers"
	. "github.com/ortin779/chirpy/models"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenExpiry  = time.Second * time.Duration(3_600)
	refreshTokenExpiry = time.Hour * 24 * 6
)

type DB struct {
	path string
	mx   *sync.RWMutex
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps       map[int]Chirp           `json:"chirps"`
	Users        map[int]User            `json:"users"`
	RefreshToken map[string]RefreshToken `json:"revoked_tokens"`
}

type NotFoundError struct{}

func (NotFoundError) Error() string {
	return "not found"
}

type AuthError struct {
	message string
}

func (aerr AuthError) Error() string {
	return aerr.message
}

func NewDB(path string) (*DB, error) {

	db := &DB{
		path: path,
		mx:   &sync.RWMutex{},
	}
	err := db.ensureDB()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbstruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	nextIndex := 1

	if len(dbstruct.Chirps) > 0 {
		keys := getSortedKeys(dbstruct.Chirps)
		nextIndex = keys[0] + 1
	}

	newChirp := Chirp{
		Id:   nextIndex,
		Body: body,
	}
	dbstruct.Chirps[nextIndex] = newChirp
	err = db.writeDB(dbstruct)
	if err != nil {
		return Chirp{}, err
	}
	return newChirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbstruct, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}

	keys := getSortedKeys(dbstruct.Chirps)
	chirps := make([]Chirp, len(keys))
	for idx, key := range keys {
		chirps[idx] = dbstruct.Chirps[key]
	}
	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbstruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbstruct.Chirps[id]
	if !ok {
		return Chirp{}, NotFoundError{}
	}
	return chirp, nil
}

func (db *DB) CreateUser(userBody UserRequestBody) (UserResponse, error) {
	dbstruct, err := db.loadDB()
	if err != nil {
		return UserResponse{}, err
	}

	existingUsr := findUser(userBody.Email, dbstruct.Users)
	if existingUsr != nil {
		return UserResponse{}, fmt.Errorf("user already exist with given email")
	}

	nextIndex := 1

	if len(dbstruct.Users) > 0 {
		keys := getSortedKeys(dbstruct.Users)
		nextIndex = keys[0] + 1
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userBody.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserResponse{}, err
	}
	newUser := User{
		Id:       nextIndex,
		Email:    userBody.Email,
		Password: string(hashedPassword),
	}
	dbstruct.Users[nextIndex] = newUser
	err = db.writeDB(dbstruct)
	if err != nil {
		return UserResponse{}, err
	}
	return UserResponse{
		Id:    newUser.Id,
		Email: newUser.Email,
	}, nil
}

func (db *DB) UpdateUser(userBody UserRequestBody, userId string) (UserResponse, error) {
	dbstruct, err := db.loadDB()
	if err != nil {
		return UserResponse{}, err
	}

	parsedId, err := strconv.Atoi(userId)
	if err != nil {
		return UserResponse{}, fmt.Errorf("invalid user id")
	}

	existingUsr, ok := dbstruct.Users[parsedId]
	if !ok {
		return UserResponse{}, fmt.Errorf("user already exist with given email")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userBody.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserResponse{}, err
	}
	updatedUser := User{
		Id:       existingUsr.Id,
		Email:    userBody.Email,
		Password: string(hashedPassword),
	}
	dbstruct.Users[parsedId] = updatedUser
	err = db.writeDB(dbstruct)
	if err != nil {
		return UserResponse{}, err
	}
	return UserResponse{
		Id:    updatedUser.Id,
		Email: updatedUser.Email,
	}, nil
}

func (db *DB) LoginUser(userBody UserRequestBody) (UserLoginResponse, error) {
	dbstruct, err := db.loadDB()
	if err != nil {
		return UserLoginResponse{}, err
	}

	user := findUser(userBody.Email, dbstruct.Users)
	if user == nil {
		return UserLoginResponse{}, AuthError{message: fmt.Sprintf("no user with given email %s", userBody.Email)}
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userBody.Password))
	if err != nil {
		return UserLoginResponse{}, AuthError{message: fmt.Sprintf("invalid password for user with email %s", userBody.Email)}
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
		return UserLoginResponse{}, errors.New("error while signing the token")
	}
	refreshToken, err := helpers.CreateToken(refreshTokenClaims)
	if err != nil {
		return UserLoginResponse{}, errors.New("error while signing the token")
	}
	dbstruct.RefreshToken[refreshToken] = RefreshToken{
		Id:         refreshToken,
		HasRevoked: false,
	}
	err = db.writeDB(dbstruct)
	if err != nil {
		return UserLoginResponse{}, err
	}
	return UserLoginResponse{
		Id:           user.Id,
		Email:        user.Email,
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (db *DB) RefreshToken(token string) (RefreshTokenResponse, error) {
	dbstruct, err := db.loadDB()
	if err != nil {
		return RefreshTokenResponse{}, err
	}

	parsedToken, err := helpers.ValidateToken(token)

	if err != nil {
		return RefreshTokenResponse{}, AuthError{message: err.Error()}
	}
	if !parsedToken.Valid {
		return RefreshTokenResponse{}, AuthError{message: "invalid refresh token"}
	}
	issuer, err := parsedToken.Claims.GetIssuer()
	if err != nil {
		return RefreshTokenResponse{}, err
	}
	if issuer != "chirpy-refresh" {
		return RefreshTokenResponse{}, AuthError{message: "invalid refresh token issuer"}
	}
	userId, err := parsedToken.Claims.GetSubject()
	if err != nil {
		return RefreshTokenResponse{}, errors.New("invalid refresh token claims")
	}

	rToken, ok := dbstruct.RefreshToken[parsedToken.Raw]
	if !ok {
		return RefreshTokenResponse{}, AuthError{message: "invalid refresh token"}
	}

	if rToken.HasRevoked {
		return RefreshTokenResponse{}, AuthError{message: "refresh token revoked"}
	}

	accessTokenClaims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExpiry)),
		Issuer:    "chirpy-access",
		Subject:   userId,
	}

	accessToken, err := helpers.CreateToken(accessTokenClaims)
	if err != nil {
		return RefreshTokenResponse{}, errors.New("error while signing the token")
	}
	return RefreshTokenResponse{Token: accessToken}, nil
}

func (db *DB) RevokeToken(token string) error {
	dbstruct, err := db.loadDB()
	if err != nil {
		return err
	}

	parsedToken, err := helpers.ValidateToken(token)

	if err != nil {
		return AuthError{message: err.Error()}
	}
	if !parsedToken.Valid {
		return AuthError{message: "invalid refresh token"}
	}
	issuer, err := parsedToken.Claims.GetIssuer()
	if err != nil {
		return err
	}
	if issuer != "chirpy-refresh" {
		return AuthError{message: "invalid refresh token issuer"}
	}

	rToken, ok := dbstruct.RefreshToken[parsedToken.Raw]
	if !ok {
		return AuthError{message: "invalid refresh token"}
	}

	if rToken.HasRevoked {
		fmt.Println("Revoked state")
		return AuthError{message: "token has been revoked"}
	}

	dbstruct.RefreshToken[parsedToken.Raw] = RefreshToken{
		Id:         rToken.Id,
		HasRevoked: true,
	}

	return db.writeDB(dbstruct)

}

func (db *DB) ensureDB() error {
	_, err := os.OpenFile(db.path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mx.Lock()
	defer db.mx.Unlock()

	file, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	if len(file) == 0 {
		return DBStructure{
			Chirps:       make(map[int]Chirp),
			Users:        make(map[int]User),
			RefreshToken: make(map[string]RefreshToken),
		}, nil
	}

	dbStructure := DBStructure{}
	err = json.Unmarshal(file, &dbStructure)
	if err != nil {
		return DBStructure{}, err
	}
	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mx.Lock()
	defer db.mx.Unlock()

	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, data, 0644)

	if err != nil {
		return err
	}

	return nil
}

func getSortedKeys[T any](m map[int]T) []int {
	keys := []int{}

	for key := range m {
		keys = append(keys, key)
	}
	slices.SortFunc(keys, func(i, j int) int { return j - i })
	return keys
}

func findUser(email string, users map[int]User) *User {
	for _, usr := range users {
		if usr.Email == email {
			return &usr
		}
	}
	return nil
}
