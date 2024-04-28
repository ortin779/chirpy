package db

import (
	"slices"
	"strconv"

	. "github.com/ortin779/chirpy/models"
)

func (db *DB) CreateChirp(body string, authorId int) (Chirp, error) {
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
		Id:       nextIndex,
		Body:     body,
		AuthorId: authorId,
	}
	dbstruct.Chirps[nextIndex] = newChirp
	err = db.writeDB(dbstruct)
	if err != nil {
		return Chirp{}, err
	}
	return newChirp, nil
}

func (db *DB) GetChirps(authorId string, sort string) ([]Chirp, error) {
	dbstruct, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}

	keys := getSortedKeys(dbstruct.Chirps)
	chirps := make([]Chirp, 0, len(keys))

	if authorId == "" {
		for _, key := range keys {
			chirps = append(chirps, dbstruct.Chirps[key])
		}
	} else {
		parsedId, err := strconv.Atoi(authorId)
		if err != nil {
			return []Chirp{}, err
		}
		for _, v := range keys {
			chirp := dbstruct.Chirps[v]
			if chirp.AuthorId == parsedId {
				chirps = append(chirps, chirp)
			}
		}
	}

	slices.SortFunc(chirps, func(a, b Chirp) int {
		if sort == "asc" {
			return a.Id - b.Id
		} else {
			return b.Id - a.Id
		}
	})

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

func (db *DB) DeleteChirp(id int, authorId int) (Chirp, error) {
	dbstruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbstruct.Chirps[id]

	if !ok {
		return Chirp{}, NotFoundError{}
	}

	if chirp.AuthorId != authorId {
		return Chirp{}, AuthorizationError{message: "you are not the author"}
	}

	delete(dbstruct.Chirps, id)

	db.writeDB(dbstruct)
	return chirp, nil
}
