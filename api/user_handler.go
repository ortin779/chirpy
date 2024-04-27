package api

import (
	"encoding/json"
	"net/http"

	"github.com/ortin779/chirpy/db"
	"github.com/ortin779/chirpy/models"
)

type UserHandler struct {
	database *db.DB
}

func NewUserHandler(db *db.DB) UserHandler {
	return UserHandler{
		database: db,
	}
}

func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	requestBody := models.UserRequestBody{}

	err := decoder.Decode(&requestBody)

	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}

	user, err := h.database.CreateUser(requestBody)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) HandleEditUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	requestBody := models.UserRequestBody{}
	userId := r.Header.Get("User-Id")

	err := decoder.Decode(&requestBody)

	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}

	user, err := h.database.UpdateUser(requestBody, userId)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, user)
}
