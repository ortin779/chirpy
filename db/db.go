package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"sync"

	"github.com/ortin779/chirpy/helpers"
	. "github.com/ortin779/chirpy/models"
	"golang.org/x/crypto/bcrypt"
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
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
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

	signedToken, err := helpers.CreateToken(userBody, strconv.Itoa(user.Id))
	if err != nil {
		fmt.Println(err)
		return UserLoginResponse{}, errors.New("error while signing the token")
	}
	return UserLoginResponse{
		Id:    user.Id,
		Email: user.Email,
		Token: signedToken,
	}, nil
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
			Chirps: make(map[int]Chirp),
			Users:  make(map[int]User),
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
