package api

import (
	"encoding/json"
	"net/http"

	"github.com/ortin779/chirpy/db"
)

type UserHandler struct {
	database *db.DB
}

type UserRequestBody struct {
	Email string `json:"email"`
}

func NewUserHandler(db *db.DB) UserHandler {
	return UserHandler{
		database: db,
	}
}

func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	requestBody := UserRequestBody{}

	err := decoder.Decode(&requestBody)

	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}

	user, err := h.database.CreateUser(requestBody.Email)
	if err != nil {
		RespondWithError(w, 500, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, user)
}
