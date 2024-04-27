package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ortin779/chirpy/db"
	"github.com/ortin779/chirpy/models"
)

type AuthHandler struct {
	database *db.DB
}

func NewAuthHandler(db *db.DB) AuthHandler {
	return AuthHandler{
		database: db,
	}
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	requestBody := models.UserRequestBody{}

	err := decoder.Decode(&requestBody)

	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}

	user, err := h.database.LoginUser(requestBody)
	if err != nil {
		if errors.As(err, &db.AuthenticationError{}) {
			RespondWithError(w, 401, err.Error())
			return
		}

		RespondWithError(w, 500, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) HandleRefresToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	if len(token) == 0 {
		RespondWithError(w, 401, "invalid refresh token")
		return
	}

	resp, err := h.database.RefreshToken(token)
	if err != nil {
		if errors.As(err, &db.AuthenticationError{}) {
			RespondWithError(w, 401, err.Error())
			return
		}

		RespondWithError(w, 500, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) HandleRevokeToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	if len(token) == 0 {
		RespondWithError(w, 401, "invalid refresh token")
		return
	}

	err := h.database.RevokeToken(token)
	if err != nil {
		if errors.As(err, &db.AuthenticationError{}) {
			RespondWithError(w, 401, err.Error())
			return
		}

		RespondWithError(w, 500, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, struct{}{})
}
