package http

import (
	"net/http"
	"server/config"

	"github.com/go-chi/chi/v5"
)

func CreateAndRunServer(r chi.Router, cfg config.HTTPConfig) error {

	httpServer := &http.Server{
		Addr:    cfg.Address,
		Handler: r,
	}

	return httpServer.ListenAndServe()
}
