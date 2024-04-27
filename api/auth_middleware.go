package api

import (
	"net/http"

	"github.com/ortin779/chirpy/helpers"
)

func AuthMiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtToken := r.Header.Get("Authorization")
		token, err := helpers.ValidateToken(jwtToken)
		if err != nil {
			RespondWithError(w, 401, err.Error())
			return
		}
		if !token.Valid {
			RespondWithError(w, 401, "invalid token")
			return
		}
		issuer, err := token.Claims.GetIssuer()
		if err != nil {
			RespondWithError(w, 401, err.Error())
			return
		}

		if issuer != "chirpy-access" {
			RespondWithError(w, 401, "invalid access token")
			return
		}
		userId, err := token.Claims.GetSubject()
		if err != nil {
			RespondWithError(w, 401, err.Error())
			return
		}

		r.Header.Set("User-Id", userId)

		next.ServeHTTP(w, r)
	})
}
