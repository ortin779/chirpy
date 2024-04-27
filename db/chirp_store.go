package db

import . "github.com/ortin779/chirpy/models"

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
