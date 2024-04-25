package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/ortin779/chirpy/db"
)

type chirpRequestBody struct {
	Body string `json:"body"`
}

var ProfaneWords = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

type ChirpHandler struct {
	database *db.DB
}

func NewChirpHandler(db *db.DB) ChirpHandler {
	return ChirpHandler{
		database: db,
	}
}

func (ch *ChirpHandler) HandleCreateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	requestBody := chirpRequestBody{}

	err := decoder.Decode(&requestBody)

	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}

	if len(requestBody.Body) > 140 {
		RespondWithError(w, 400, "Chirp is too long")
		return
	}

	chirp, err := ch.database.CreateChirp(requestBody.Body)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, chirp)
}

func (ch *ChirpHandler) HandleGetChirps(w http.ResponseWriter, r *http.Request) {

	chirps, err := ch.database.GetChirps()
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, chirps)
}

func (ch *ChirpHandler) HandleGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpId")
	parsedId, err := strconv.Atoi(chirpId)
	if err != nil {
		RespondWithError(w, 400, err.Error())
	}
	chirp, err := ch.database.GetChirp(parsedId)
	if err != nil {
		if errors.Is(err, db.NotFoundError{}) {

			RespondWithError(w, 404, err.Error())
		} else {
			RespondWithError(w, 500, err.Error())
		}
		return
	}

	RespondWithJSON(w, http.StatusOK, chirp)
}
