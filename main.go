package main

import (
	"fmt"
	"net/http"

	"github.com/ortin779/chirpy/api"
	"github.com/ortin779/chirpy/app"
)

func main() {
	apiCfg := app.ApiConfig{}
	mux := http.NewServeMux()
	corsMux := api.MiddlewareCors(mux)

	mux.Handle("/app/*", apiCfg.MiddlewareMetricInc(app.HandleFileServer()))
	mux.HandleFunc("GET /api/healthz", api.HealthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.MetricsHandler)
	mux.HandleFunc("/api/reset", apiCfg.ResetHandler)
	mux.HandleFunc("POST /api/validate_chirp", api.ValidateChirp)

	fmt.Println("Starting server on 8080")
	http.ListenAndServe(":8080", corsMux)
}
