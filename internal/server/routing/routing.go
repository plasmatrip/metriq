package routing

import (
	"github.com/go-chi/chi/v5"
	"github.com/plasmatrip/metriq/internal/logger"
	"github.com/plasmatrip/metriq/internal/server/compress"
	"github.com/plasmatrip/metriq/internal/server/config"
	"github.com/plasmatrip/metriq/internal/server/handlers"
	"github.com/plasmatrip/metriq/internal/storage"
)

func NewRouter(l *logger.Logger, s storage.Repository, c config.Config) *chi.Mux {

	h := handlers.NewHandlers(s, c)

	r := chi.NewRouter()

	r.Use(compress.WithCompressed)
	r.Use(l.WithLogging)

	r.Route("/update", func(r chi.Router) {
		r.Post("/", h.JSONUpdateHandler)
	})
	r.Post("/update/{metricType}/{metricName}/{metricValue}", h.UpdateHandler)
	r.Route("/value", func(r chi.Router) {
		r.Post("/", h.JSONValueHandler)
	})
	r.Get("/value/{metricType}/{metricName}", h.ValueHandler)
	r.Get("/", h.MetricsHandler)

	return r
}
