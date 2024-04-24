package app

import (
	"fmt"
	"html/template"
	"net/http"
)

type ApiConfig struct {
	FileServerHits int
}

func (cfg *ApiConfig) MiddlewareMetricInc(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileServerHits++
		h.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	template, err := template.ParseFiles("./app/admin/metrics.gohtml")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	err = template.Execute(w, cfg)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
}

func (cfg *ApiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/palin; charset=utf-8")
	cfg.FileServerHits = 0
	w.Write([]byte("Hits: " + fmt.Sprint(cfg.FileServerHits)))
}

func HandleFileServer() http.Handler {
	return http.FileServer(http.Dir("."))
}
