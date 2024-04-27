package db

import (
	"encoding/json"
	"os"
	"slices"
	"sync"
	"time"

	. "github.com/ortin779/chirpy/models"
)

const (
	accessTokenExpiry  = time.Second * time.Duration(3_600)
	refreshTokenExpiry = time.Hour * 24 * 6
)

type DB struct {
	path string
	mx   *sync.RWMutex
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

type AuthenticationError struct {
	message string
}

type AuthorizationError struct {
	message string
}

func (authErr AuthenticationError) Error() string {
	return authErr.message
}

func (aerr AuthorizationError) Error() string {
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
