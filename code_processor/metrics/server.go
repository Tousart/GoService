package metrics

import (
	"fmt"
	"httpServer/code_processor/config"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Listen(cfg config.Prometheus) error {
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	r := chi.NewRouter()

	r.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting metrics server on %s", address)
	return http.ListenAndServe(address, r)
}
