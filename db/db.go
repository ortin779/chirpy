package db

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"sync"
)

type DB struct {
	path string
	mx   *sync.RWMutex
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

type NotFoundError struct{}

func (NotFoundError) Error() string {
	return "not found"
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

	fmt.Println("Calling create chirp")

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

func (db *DB) CreateUser(email string) (User, error) {
	dbstruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	nextIndex := 1

	if len(dbstruct.Users) > 0 {
		keys := getSortedKeys(dbstruct.Users)
		nextIndex = keys[0] + 1
	}

	newUser := User{
		Id:    nextIndex,
		Email: email,
	}
	dbstruct.Users[nextIndex] = newUser
	err = db.writeDB(dbstruct)
	if err != nil {
		return User{}, err
	}
	return newUser, nil
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
		fmt.Println(err)
		return DBStructure{}, err
	}
	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mx.Lock()
	defer db.mx.Unlock()

	data, err := json.Marshal(dbStructure)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = os.WriteFile(db.path, data, 0644)

	if err != nil {
		fmt.Println(err.Error())
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
