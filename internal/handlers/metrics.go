package handlers

import (
	"httpServer/internal/mux1"
	"net/http"

	"github.com/VictoriaMetrics/metrics"
)

func RegisterMetricsHandler(mux *mux1.ServeMux) {

	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics.WritePrometheus(w, true)
	})

}
