package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/ortin779/chirpy/api"
	"github.com/ortin779/chirpy/app"
	"github.com/ortin779/chirpy/db"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("error while loading env variables")
	}
	apiCfg := app.ApiConfig{}
	mux := http.NewServeMux()
	corsMux := api.MiddlewareCors(mux)
	database, err := db.NewDB("database.json")
	if err != nil {
		log.Fatalf(err.Error())
	}

	chirpHandler := api.NewChirpHandler(database)
	userHandler := api.NewUserHandler(database)
	authHandler := api.NewAuthHandler(database)
	polkaHanler := api.NewPolksHandler(database)

	mux.Handle("/app/*", apiCfg.MiddlewareMetricInc(app.HandleFileServer()))
	mux.HandleFunc("GET /api/healthz", api.HealthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.MetricsHandler)
	mux.HandleFunc("/api/reset", apiCfg.ResetHandler)

	mux.Handle("POST /api/chirps", api.AuthMiddleware(chirpHandler.HandleCreateChirp))
	mux.HandleFunc("GET /api/chirps", chirpHandler.HandleGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpId}", chirpHandler.HandleGetChirp)
	mux.Handle("DELETE /api/chirps/{chirpId}", api.AuthMiddleware(chirpHandler.HandleDeleteChirp))

	mux.HandleFunc("POST /api/users", userHandler.HandleCreateUser)
	mux.Handle("PUT /api/users", api.AuthMiddleware(userHandler.HandleEditUser))

	mux.HandleFunc("POST /api/login", authHandler.HandleLogin)
	mux.HandleFunc("POST /api/refresh", authHandler.HandleRefresToken)
	mux.HandleFunc("POST /api/revoke", authHandler.HandleRevokeToken)

	mux.HandleFunc("POST /api/polka/webhooks", polkaHanler.HandlePolkaWebhook)

	fmt.Println("Starting server on 8080")
	http.ListenAndServe(":8080", corsMux)
}
