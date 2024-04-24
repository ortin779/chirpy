package api

import "net/http"

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/palin; charset=utf-8")
	w.Write([]byte("OK"))
}
