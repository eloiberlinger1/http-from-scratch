package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/VictoriaMetrics/metrics"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // statut par défaut
		}

		// Exécute le handler réel
		next.ServeHTTP(rec, r)

		// Formate les labels dynamiquement (méthode, chemin, code statut)
		labels := fmt.Sprintf(`path="%s",method="%s",status="%d"`, r.URL.Path, r.Method, rec.statusCode)

		// 1. Incrémente le compteur total de requêtes
		metrics.GetOrCreateCounter(fmt.Sprintf(`http_requests_total{%s}`, labels)).Inc()

		// 2. Met à jour l'histogramme de durée d'exécution (en secondes)
		metrics.GetOrCreateHistogram(fmt.Sprintf(`http_request_duration_seconds{%s}`, labels)).UpdateDuration(start)
	})
}
