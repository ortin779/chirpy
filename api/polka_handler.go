package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ortin779/chirpy/db"
	"github.com/ortin779/chirpy/models"
)

type PolkaHandler struct {
	database *db.DB
}

func NewPolksHandler(db *db.DB) PolkaHandler {
	return PolkaHandler{
		database: db,
	}
}

func (ph PolkaHandler) HandlePolkaWebhook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var polkaBody models.PolkaBody

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&polkaBody)

	if err != nil {
		RespondWithError(w, 400, "invalid polka body")
		return
	}

	if polkaBody.Event != "user.upgraded" {
		RespondWithJSON(w, 200, struct{}{})
		return
	}

	err = ph.database.MarkUserAsRedChirp(polkaBody.Data.UserId)

	if err != nil {
		if errors.As(err, &db.NotFoundError{}) {
			RespondWithError(w, 404, err.Error())
			return
		}
		RespondWithError(w, 500, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, struct{}{})
}
