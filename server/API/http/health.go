package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func checkHealthy(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func Health(r chi.Router) {
	r.Get("/health", checkHealthy)
}
