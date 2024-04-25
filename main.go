package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ortin779/chirpy/api"
	"github.com/ortin779/chirpy/app"
	"github.com/ortin779/chirpy/db"
)

func main() {
	apiCfg := app.ApiConfig{}
	mux := http.NewServeMux()
	corsMux := api.MiddlewareCors(mux)
	database, err := db.NewDB("database.json")
	if err != nil {
		log.Fatalf(err.Error())
	}

	chirpHandler := api.NewChirpHandler(database)

	mux.Handle("/app/*", apiCfg.MiddlewareMetricInc(app.HandleFileServer()))
	mux.HandleFunc("GET /api/healthz", api.HealthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.MetricsHandler)
	mux.HandleFunc("/api/reset", apiCfg.ResetHandler)
	mux.HandleFunc("POST /api/chirps", chirpHandler.HandleCreateChirp)
	mux.HandleFunc("GET /api/chirps", chirpHandler.HandleGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpId}", chirpHandler.HandleGetChirp)

	fmt.Println("Starting server on 8080")
	http.ListenAndServe(":8080", corsMux)
}
