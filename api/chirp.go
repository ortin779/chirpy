package api

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

type chirpRequestBody struct {
	Body string `json:"body"`
}

var ProfaneWords = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

func ValidateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	requestBody := chirpRequestBody{}

	err := decoder.Decode(&requestBody)

	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}

	strParts := strings.Split(requestBody.Body, " ")
	for _, word := range ProfaneWords {
		ind := slices.IndexFunc(strParts, func(v string) bool {
			return strings.ToLower(v) == word
		})
		if ind != -1 {
			strParts[ind] = "****"
		}
	}

	cleanedString := strings.Join(strParts, " ")

	if len(requestBody.Body) > 140 {
		RespondWithError(w, 400, "Chirp is too long")
		return
	}

	RespondWithJSON(w, http.StatusOK, struct {
		CleanedBody string `json:"cleaned_body"`
	}{CleanedBody: cleanedString})
}
