package service

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "code_processor",
		Name:      "request duration",
		Buckets:   prometheus.DefBuckets, // Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}
	}, []string{"translator"})
	requestTranslator = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "code_processor",
		Name:      "used translator",
	}, []string{"translator"})
)

func observeReqDuration(duration time.Duration, translator string) {
	requestDuration.WithLabelValues(translator).Observe(duration.Seconds())
}

func observeUsedTranslator(translator string) {
	requestTranslator.WithLabelValues(translator).Inc()
}
